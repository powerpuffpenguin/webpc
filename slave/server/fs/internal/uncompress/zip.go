package uncompress

import (
	"archive/zip"
	"io"
	"os"
)

type ZipReader struct {
	r *zip.Reader
}

func NewZipReader(f *os.File) (*ZipReader, error) {
	ret, e := f.Seek(0, os.SEEK_END)
	if e != nil {
		return nil, e
	}
	zr, e := zip.NewReader(f, ret)
	if e != nil {
		return nil, e
	}

	return &ZipReader{
		r: zr,
	}, nil
}
func (zr *ZipReader) Root(callback func(io.Reader, string, os.FileMode) error) (e error) {
	var r io.ReadCloser
	for _, zipFile := range zr.r.File {
		name := zipFile.Name
		perm := zipFile.Mode()

		if perm.IsDir() {
			e = callback(nil, name, perm)
			if e != nil {
				break
			}
		} else {
			r, e = zipFile.Open()
			if e != nil {
				break
			}
			e = callback(r, name, perm)
			r.Close()
			if e != nil {
				break
			}
		}
	}
	return
}
