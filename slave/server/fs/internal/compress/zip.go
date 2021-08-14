package compress

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/single/mount"
)

type Compressor interface {
	Root(dir string, name string) (e error)
	Close() error
}
type ZipWriter struct {
	helper
	m *mount.Mount
	w *zip.Writer
}

func NewZipWriter(helper helper, m *mount.Mount, w io.Writer) *ZipWriter {
	return &ZipWriter{
		helper: helper,
		m:      m,
		w:      zip.NewWriter(w),
	}
}
func (zw *ZipWriter) Close() error {
	return zw.w.Close()
}

func (zw *ZipWriter) Root(vdir string, name string) (e error) {
	m := zw.m
	dir, e := m.Filename(vdir)
	if e != nil {
		return
	}
	root := filepath.Join(dir, name)
	count := len(dir)
	e = m.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := path[count:]
		err = zw.SendProgress(name)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = zw.dir(info, name)
		} else {
			err = zw.file(info, path, name)
		}
		return err
	})
	return
}
func (zw *ZipWriter) dir(info os.FileInfo, name string) (e error) {
	header, e := zip.FileInfoHeader(info)
	if e != nil {
		return e
	}
	header.Name = name
	_, e = zw.w.CreateHeader(header)
	return e
}

func (zw *ZipWriter) file(info os.FileInfo, filename, name string) (e error) {
	m := zw.m
	f, e := m.Open(filename)
	if e != nil {
		return
	}
	defer f.Close()

	return e
}
