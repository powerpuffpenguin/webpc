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

func (zw *ZipWriter) Root(dir string, name string) (e error) {
	m := zw.m
	root := filepath.Join(dir, name)
	count := len(dir)
	var f *os.File
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
			f, err = m.Open(path)
			if err != nil {
				return err
			}
			err = zw.file(info, f, name)
			f.Close()
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

func (zw *ZipWriter) file(info os.FileInfo, reader io.Reader, name string) (e error) {
	header, e := zip.FileInfoHeader(info)
	if e != nil {
		return
	}
	header.Name = name
	w, e := zw.w.CreateHeader(header)
	if e != nil {
		return
	}
	_, e = io.Copy(w, reader)
	return e
}
