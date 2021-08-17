package copied

import (
	"io"
	"os"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (w *Worker) copyFile(m *mount.Mount, dst string, perm os.FileMode, r io.Reader) (e error) {
	f, e := w.openFile(m, dst, perm)
	if e != nil {
		return
	} else if f == nil {
		return
	}
	_, e = io.Copy(f, r)
	f.Close()
	return
}
func (w *Worker) openFile(m *mount.Mount, name string, perm os.FileMode) (f *os.File, e error) {
	if w.style == grpc_fs.Event_YesAll {
		f, e = m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		return
	}
	f, e = m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, perm)
	if status.Code(e) != codes.AlreadyExists {
		return
	}
	if w.style == grpc_fs.Event_SkipAll {
		e = nil
		return
	}

	// send exists
	e = w.server.Send(&grpc_fs.CopyResponse{
		Event: grpc_fs.Event_Exists,
		Value: name,
	})
	if e != nil {
		return
	}

	// wait choice
	req, e := w.waitRequest(grpc_fs.Event_Yes, grpc_fs.Event_YesAll,
		grpc_fs.Event_Skip, grpc_fs.Event_SkipAll,
		grpc_fs.Event_No,
	)
	if e != nil {
		return
	}
	w.style = req.Event

	if req.Event == grpc_fs.Event_No {
		e = status.Error(codes.Canceled, name+` already exists, cancel uncompress`)
		return
	} else if req.Event == grpc_fs.Event_Skip || req.Event == grpc_fs.Event_SkipAll {
		return
	}

	f, e = m.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	return
}
