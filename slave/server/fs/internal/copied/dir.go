package copied

import (
	"os"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (w *Worker) copyDir(m *mount.Mount, dst string, perm os.FileMode) (e error) {
	info, e := m.Stat(dst)
	if e == nil { //exists
		if info.IsDir() {
			if info.Mode() != perm {
				e = w.onExists(m, dst, perm)
			}
		} else {
			e = status.Error(codes.AlreadyExists, `file already exists: `+dst)
		}
		return
	}
	if codes.NotFound != status.Code(e) {
		return
	}
	path, e := m.Filename(dst)
	if e != nil {
		return
	}
	e = os.Mkdir(path, perm)
	if e != nil {
		if os.IsNotExist(e) {
			e = status.Error(codes.NotFound, `not exists: `+dst)
		} else if os.IsExist(e) {
			e = status.Error(codes.AlreadyExists, `already exists: `+dst)
		} else if os.IsPermission(e) {
			e = status.Error(codes.PermissionDenied, `forbidden: `+dst)
		}
	}
	return
}
func (w *Worker) onExists(m *mount.Mount, name string, perm os.FileMode) (e error) {
	if w.style == grpc_fs.Event_YesAll {
		e = m.Chmod(name, perm)
		return
	} else if w.style == grpc_fs.Event_SkipAll {
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
		e = status.Error(codes.Canceled, name+` already exists, cancel copy`)
		return
	} else if req.Event == grpc_fs.Event_Skip || req.Event == grpc_fs.Event_SkipAll {
		return
	}
	e = m.Chmod(name, perm)
	return
}
