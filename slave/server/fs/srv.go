package fs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/sessions"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var modtime = time.Now()

type server struct {
	grpc_fs.UnimplementedFSServer
	helper.Helper
}

var emptyMountResponse grpc_fs.MountResponse

func (s server) checkRead(ctx context.Context, m *mount.Mount) (e error) {
	// shared
	if m.Shared() {
		return
	}

	// can read
	if !m.Read() {
		e = s.Error(codes.PermissionDenied, `filesystem is not readable`)
		return
	}

	// root
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	if userdata.AuthTest(db.Root) {
		return
	}

	// read auth
	if !userdata.AuthTest(db.Read) {
		e = s.Error(codes.PermissionDenied, `no read permission`)
		return
	}
	return
}
func (s server) checkWrite(ctx context.Context, m *mount.Mount) (e error) {
	// can read
	if !m.Write() {
		e = s.Error(codes.PermissionDenied, `filesystem is not writable`)
		return
	}

	// root
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	if userdata.AuthTest(db.Root) {
		return
	}

	// read auth
	if !userdata.AuthTest(db.Write) {
		e = s.Error(codes.PermissionDenied, `no read permission`)
		return
	}
	return
}
func (s server) checkUserdataWrite(userdata *sessions.Userdata, m *mount.Mount) (e error) {
	// can read
	if !m.Write() {
		e = s.Error(codes.PermissionDenied, `filesystem is not writable`)
		return
	}

	if userdata.AuthTest(db.Root) {
		return
	}

	// read auth
	if !userdata.AuthTest(db.Write) {
		e = s.Error(codes.PermissionDenied, `no read permission`)
		return
	}
	return
}
func (s server) Mount(ctx context.Context, req *grpc_fs.MountRequest) (resp *grpc_fs.MountResponse, e error) {
	fs := mount.Default()

	s.SetHTTPCacheMaxAge(ctx, 5)
	e = s.ServeMessage(ctx, modtime, func(nobody bool) error {
		if nobody {
			resp = &emptyMountResponse
		} else {
			resp = &grpc_fs.MountResponse{
				Name: fs.Names(),
			}
		}
		return nil
	})
	return
}

var emptyListResponse grpc_fs.ListResponse

func (s server) List(ctx context.Context, req *grpc_fs.ListRequest) (resp *grpc_fs.ListResponse, e error) {
	fs := mount.Default()
	m := fs.Root(req.Root)
	if m == nil {
		e = s.Error(codes.NotFound, `root not found: `+req.Root)
		return
	}
	e = s.checkRead(ctx, m)
	if e != nil {
		return
	}

	dir, modtime, items, e := m.LS(req.Path)
	if e != nil {
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 5)
	e = s.ServeMessage(ctx, modtime, func(nobody bool) error {
		if nobody {
			resp = &emptyListResponse
		} else {
			resp = &grpc_fs.ListResponse{
				Dir: &grpc_fs.Dir{
					Root:   m.Name(),
					Read:   m.Read(),
					Write:  m.Write(),
					Shared: m.Shared(),
					Dir:    dir,
				},
			}
			if len(items) != 0 {
				resp.Items = make([]*grpc_fs.FileInfo, len(items))
				for i, item := range items {
					resp.Items[i] = &grpc_fs.FileInfo{
						Name:  item.Name,
						Mode:  item.Mode,
						Size:  item.Size,
						IsDir: item.IsDir,
					}
				}
			}
		}
		return nil
	})
	return
}

func (s server) Download(req *grpc_fs.DownloadRequest, server grpc_fs.FS_DownloadServer) (e error) {
	fs := mount.Default()
	m := fs.Root(req.Root)
	if m == nil {
		e = s.Error(codes.NotFound, `root not found: `+req.Root)
		return
	}
	ctx := server.Context()
	e = s.checkRead(ctx, m)
	if e != nil {
		return
	}
	f, e := m.Open(req.Path)
	if e != nil {
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		e = s.ToHTTPError(req.Path, e)
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 0)

	grpc.SetHeader(ctx, metadata.Pairs(
		`Content-Disposition`,
		fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(req.Path)),
	))
	e = s.ServeName(server,
		req.Path,
		stat.ModTime(),
		f,
	)
	return
}
func (s server) putFirst(userdata *sessions.Userdata, root, path string) (f *os.File, e error) {
	fs := mount.Default()
	m := fs.Root(root)
	if m == nil {
		e = s.Error(codes.NotFound, `root not found: `+root)
		return
	}
	e = s.checkUserdataWrite(userdata, m)
	if e != nil {
		return
	}
	f, e = m.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0666)
	if e != nil {
		return
	}
	return
}
func (s server) closeFile(f *os.File, e error) error {
	if f == nil {
		return e
	}

	if e == nil {
		es := f.Sync()
		ec := f.Close()

		if es == nil {
			e = ec
		} else {
			e = es
		}
	} else {
		e = f.Close()
	}
	return e
}
func (s server) Put(server grpc_fs.FS_PutServer) (e error) {
	TAG := `forward.fs Put`
	_, userdata, e := s.JSONUserdata(server.Context())
	if e != nil {
		return
	}

	var (
		first  = true
		req    grpc_fs.PutRequest
		f      *os.File
		root   string
		path   string
		writed uint64
		n      int
	)
	for {
		e = server.RecvMsg(&req)
		if e != nil {
			if e == io.EOF {
				e = nil
			}
			break
		}
		if first {
			first = false
			root = req.Root
			path = req.Path
			f, e = s.putFirst(&userdata, req.Root, req.Path)
			if e != nil {
				break
			}
			if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
				ce.Write(
					zap.String(`event`, `open`),
					zap.String(`who`, userdata.Who()),
					zap.String(`root`, req.Root),
					zap.String(`path`, req.Path),
				)
			}
		}
		if len(req.Data) != 0 {
			n, e = f.Write(req.Data)
			writed += uint64(n)
			if e != nil {
				break
			}
		}
	}
	e = s.closeFile(f, e)
	if e == nil {
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`event`, `success`),
				zap.String(`who`, userdata.Who()),
				zap.String(`root`, root),
				zap.String(`path`, path),
				zap.Uint64(`writed`, writed),
			)
		}
	} else {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`event`, `error`),
				zap.String(`who`, userdata.Who()),
				zap.String(`root`, root),
				zap.String(`path`, path),
				zap.Uint64(`writed`, writed),
			)
		}
	}
	return
}
