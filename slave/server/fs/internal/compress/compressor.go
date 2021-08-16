package compress

import (
	"io"
	"os"
	"path/filepath"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
)

type writer interface {
	dir(info os.FileInfo, name string) (e error)
	file(info os.FileInfo, reader io.Reader, name string) (e error)
	Close() (e error)
}
type Compressor struct {
	server grpc_fs.FS_CompressServer
	w      writer
	m      *mount.Mount
}

func NewCompressor(server grpc_fs.FS_CompressServer, m *mount.Mount, w writer) *Compressor {
	return &Compressor{
		server: server,
		w:      w,
		m:      m,
	}
}
func (c *Compressor) Close() error {
	return c.w.Close()
}
func (c *Compressor) Root(dir string, name string) (e error) {
	m := c.m
	root := filepath.Join(dir, name)
	count := len(dir)
	var f *os.File
	e = m.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := path[count:]
		err = c.server.Send(&grpc_fs.CompressResponse{
			Event: grpc_fs.Event_Progress,
			Value: name,
		})
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = c.w.dir(info, name)
		} else {
			f, err = m.Open(path)
			if err != nil {
				return err
			}
			err = c.w.file(info, f, name)
			f.Close()
		}
		return err
	})
	return
}
