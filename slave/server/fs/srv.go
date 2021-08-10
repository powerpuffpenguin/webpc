package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
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
	// root
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	if userdata.AuthTest(db.Root) {
		return
	}

	// can read
	if !m.Read() {
		e = s.Error(codes.PermissionDenied, `filesystem is not readable`)
		return
	}

	// read auth
	if !userdata.AuthTest(db.Read) {
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
	filename, e := m.Filename(req.Path)
	if e != nil {
		return
	}
	f, e := os.Open(filename)
	if e != nil {
		e = s.ToHTTPError(ctx, req.Path, e)
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		e = s.ToHTTPError(ctx, req.Path, e)
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 5)

	grpc.SetHeader(ctx, metadata.Pairs(
		`Content-Disposition`,
		fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(req.Path)),
	))
	e = s.ServeContent(server,
		`application/octet-stream`,
		stat.ModTime(),
		f,
	)
	return
}
