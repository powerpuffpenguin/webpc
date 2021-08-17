package copied

import (
	"io/fs"
	"os"
	"path/filepath"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
)

func (w *Worker) serve(mw, mr *mount.Mount) error {
	if w.Copied {
		return w.copy(mw, mr)
	} else {
		return w.move(mw, mr)
	}
}

func (w *Worker) copy(mw, mr *mount.Mount) (e error) {
	var (
		src  string
		info os.FileInfo
		f    *os.File
	)
	for _, name := range w.Names {
		src = filepath.Join(w.SrcDir, name)
		info, e = mr.Stat(src)
		if e != nil {
			break
		}
		if info.IsDir() {
			e = w.rootDir(mw, mr, w.SrcDir, name)
		} else {
			e = w.server.Send(&grpc_fs.CopyResponse{
				Event: grpc_fs.Event_Progress,
				Value: mount.Separator + name,
			})
			if e != nil {
				break
			}
			f, e = mr.Open(src)
			if e != nil {
				break
			}
			e = w.copyFile(mw, filepath.Join(w.DstDir, name), info.Mode(), f)
			f.Close()
		}
		if e != nil {
			break
		}
	}
	return
}
func (w *Worker) rootDir(mw, mr *mount.Mount, dir, name string) error {
	root := filepath.Join(dir, name)
	count := len(dir)
	return mr.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		name := path[count:]
		err = w.server.Send(&grpc_fs.CopyResponse{
			Event: grpc_fs.Event_Progress,
			Value: name,
		})
		if err != nil {
			return err
		}
		dst := filepath.Join(w.DstDir, name)
		if info.IsDir() {
			return w.copyDir(mw, dst, info.Mode())
		} else {
			f, err := mr.Open(path)
			if err != nil {
				return err
			}
			err = w.copyFile(mw, dst, info.Mode(), f)
			f.Close()
			return err
		}
	})
}
