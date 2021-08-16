package uncompress

import (
	"archive/tar"
	"io"
	"os"
)

type TarReader struct {
	r *tar.Reader
}

func NewTarReader(r io.Reader) (*TarReader, error) {
	return &TarReader{
		r: tar.NewReader(r),
	}, nil
}
func (tr *TarReader) Root(callback func(io.Reader, string, os.FileMode) error) (e error) {
	var (
		r      = tr.r
		header *tar.Header
	)
	for {
		header, e = r.Next()
		if e != nil {
			if e == io.EOF {
				e = nil
			}
			break
		}
		switch header.Typeflag {
		case tar.TypeDir:
			e = callback(nil, header.Name, os.FileMode(header.Mode))
			if e != nil {
				return
			}
		case tar.TypeReg:
			e = callback(r, header.Name, os.FileMode(header.Mode))
			if e != nil {
				return
			}
		}
	}
	return
}
