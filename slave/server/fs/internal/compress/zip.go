package compress

import (
	"archive/zip"
	"io"
	"os"
)

type ZipWriter struct {
	w *zip.Writer
}

func NewZipWriter(w io.Writer) *ZipWriter {
	return &ZipWriter{
		w: zip.NewWriter(w),
	}
}
func (zw *ZipWriter) Close() error {
	return zw.w.Close()
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
