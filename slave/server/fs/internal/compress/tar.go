package compress

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

type TarWriter struct {
	w  *tar.Writer
	gz *gzip.Writer
}

func NewTarWriter(w io.Writer, gz bool) *TarWriter {
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
		w:  writer,
		gz: gwriter,
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
