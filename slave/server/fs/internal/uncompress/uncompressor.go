package uncompress

import (
	"io"
	"os"
	"path/filepath"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type reader interface {
	Root(func(io.Reader, string, os.FileMode) error) error
}
type Uncompressor struct {
	w     *Worker
	r     reader
	m     *mount.Mount
	style grpc_fs.Event
}

func NewUncompressor(w *Worker, m *mount.Mount, r reader) *Uncompressor {
	return &Uncompressor{
		w:     w,
		r:     r,
		m:     m,
		style: grpc_fs.Event_EventUniversal,
	}
}
func (un *Uncompressor) Root(dir string) (e error) {
	r := un.r
	e = r.Root(un.root)
	return
}
func (un *Uncompressor) root(r io.Reader, path string, perm os.FileMode) (e error) {
	e = un.w.server.Send(&grpc_fs.UncompressResponse{
		Event: grpc_fs.Event_Progress,
		Value: path,
	})
	if e != nil {
		return
	}
	if r == nil {
		path = filepath.Join(un.w.Dir, path)
		e = un.m.SyncDir(path, perm)
	} else {
		e = un.file(path, r, perm)
	}
	return
}
func (un *Uncompressor) file(path string, r io.Reader, perm os.FileMode) (e error) {
	f, e := un.openFile(path, perm)
	if e != nil {
		return
	} else if f == nil {
		return
	}
	_, e = io.Copy(f, r)
	f.Close()
	return
}
func (un *Uncompressor) openFile(path string, perm os.FileMode) (f *os.File, e error) {
	name := filepath.Join(un.w.Dir, path)
	if un.style == grpc_fs.Event_YesAll {
		f, e = un.m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		return
	}
	f, e = un.m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, perm)
	if status.Code(e) != codes.AlreadyExists {
		return
	}
	if un.style == grpc_fs.Event_SkipAll {
		e = nil
		return
	}

	// send exists
	e = un.w.server.Send(&grpc_fs.UncompressResponse{
		Event: grpc_fs.Event_Exists,
		Value: path,
	})
	if e != nil {
		return
	}

	// wait choice
	req, e := un.w.waitRequest(grpc_fs.Event_Yes, grpc_fs.Event_YesAll,
		grpc_fs.Event_Skip, grpc_fs.Event_SkipAll,
		grpc_fs.Event_No,
	)
	if e != nil {
		return
	}
	un.style = req.Event

	if req.Event == grpc_fs.Event_No {
		e = status.Error(codes.Canceled, name+` already exists, cancel uncompress`)
		return
	} else if req.Event == grpc_fs.Event_Skip || req.Event == grpc_fs.Event_SkipAll {
		return
	}

	f, e = un.m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	return
}
