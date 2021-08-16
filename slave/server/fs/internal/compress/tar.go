package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/single/mount"
)

type TarWriter struct {
	helper
	m  *mount.Mount
	w  *tar.Writer
	gz *gzip.Writer
}

func NewTarWriter(helper helper, m *mount.Mount, w io.Writer, gz bool) *TarWriter {
	var (
		gwriter *gzip.Writer
		writer  *tar.Writer
	)
	if gz {
		gwriter = gzip.NewWriter(w)
		writer = tar.NewWriter(gwriter)
	} else {
		writer = tar.NewWriter(w)
	}
	return &TarWriter{
		helper: helper,
		m:      m,
		w:      writer,
		gz:     gwriter,
	}
}
func (tw *TarWriter) Close() (e error) {
	e = tw.w.Close()
	if tw.gz != nil {
		egz := tw.gz.Close()
		if e == nil {
			e = egz
		}
	}
	return
}
func (tw *TarWriter) Root(dir string, name string) (e error) {
	m := tw.m
	root := filepath.Join(dir, name)
	count := len(dir)
	var f *os.File
	e = m.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := path[count:]
		err = tw.SendProgress(name)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = tw.dir(info, name)
		} else {
			f, err = m.Open(path)
			if err != nil {
				return err
			}
			err = tw.file(info, f, name)
			f.Close()
		}
		return err
	})
	return
}
func (tw *TarWriter) dir(info os.FileInfo, name string) (e error) {
	header := &tar.Header{
		Typeflag: tar.TypeDir,
		Name:     name,
		Mode:     int64(info.Mode()),
		Uid:      os.Getuid(),
		Gid:      os.Getgid(),
		Size:     info.Size(),
		ModTime:  info.ModTime(),
	}
	return tw.w.WriteHeader(header)
}
func (tw *TarWriter) file(info os.FileInfo, reader io.Reader, name string) (e error) {
	header := &tar.Header{
		Name:    name,
		Mode:    int64(info.Mode()),
		Uid:     os.Getuid(),
		Gid:     os.Getgid(),
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}
	if e = tw.w.WriteHeader(header); e != nil {
		return
	}
	_, e = io.Copy(tw.w, reader)
	return
}
