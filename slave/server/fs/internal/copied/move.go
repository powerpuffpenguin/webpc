package copied

import (
	"os"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (w *Worker) move(mw, mr *mount.Mount) (e error) {
	var (
		dst, src, osdst, ossrc string
	)
	for _, name := range w.Names {
		dst = filepath.Join(w.DstDir, name)
		osdst, e = mw.Filename(dst)
		if e != nil {
			break
		}
		src = filepath.Join(w.SrcDir, name)
		ossrc, e = mr.Filename(src)
		if e != nil {
			break
		}

		_, err := os.Stat(osdst)
		if !os.IsNotExist(err) {
			if os.IsPermission(err) {
				e = status.Error(codes.PermissionDenied, `forbidden: `+dst)
			} else {
				e = status.Error(codes.AlreadyExists, `already exists: `+dst)
			}
			return
		}
		e = os.Rename(ossrc, osdst)
		if e != nil {
			if os.IsNotExist(e) {
				e = status.Error(codes.NotFound, `not exists: `+src)
			} else if os.IsExist(e) {
				e = status.Error(codes.AlreadyExists, `already exists: `+dst)
			} else if os.IsPermission(e) {
				e = status.Error(codes.PermissionDenied, `forbidden move: `+src+` -> `+dst)
			}
			break
		}
	}
	return
}
