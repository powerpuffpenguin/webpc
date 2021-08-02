package web

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc/codes"
	"gopkg.in/yaml.v2"
)

func (h Helper) ToHTTPError(c *gin.Context, name string, e error) {
	if os.IsNotExist(e) {
		h.Error(c, http.StatusNotFound, codes.NotFound, `not exists : `+name)
		return
	}
	if os.IsExist(e) {
		h.Error(c, http.StatusForbidden, codes.PermissionDenied, `already exists : `+name)
		return
	}
	if os.IsPermission(e) {
		h.Error(c, http.StatusForbidden, codes.PermissionDenied, `forbidden : `+name)
		return
	}
	h.Error(c, http.StatusInternalServerError, codes.Unknown, e.Error())
}

func (h Helper) NegotiateFilesystem(c *gin.Context, fs http.FileSystem, path string, index bool) {
	if path == `/` || path == `` {
		path = `/index.html`
	}
	f, e := fs.Open(path)
	if e != nil {
		if !index {
			h.ToHTTPError(c, path, e)
			return
		}
		if path != `/index.html` && os.IsNotExist(e) {
			path = `/index.html`
			f, e = fs.Open(path)
		}
	}
	if e != nil {
		h.ToHTTPError(c, path, e)
		return
	}
	stat, e := f.Stat()
	if e != nil {
		f.Close()
		h.ToHTTPError(c, path, e)
		return
	}
	if stat.IsDir() {
		f.Close()
		h.Error(c, http.StatusForbidden, codes.PermissionDenied, `not a file`)
		return
	}

	_, name := filepath.Split(path)
	http.ServeContent(c.Writer, c.Request, name, stat.ModTime(), f)
	f.Close()
}

func (h Helper) NegotiateObject(c *gin.Context, modtime time.Time, name string, obj interface{}) {
	reader := &objectReader{
		obj: obj,
	}
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEXML:
		c.Writer.Header().Set(`Content-Type`, `application/xml; charset=utf-8`)
		reader.marshal = xml.Marshal
	case binding.MIMEYAML:
		c.Writer.Header().Set(`Content-Type`, `application/x-yaml; charset=utf-8`)
		reader.marshal = yaml.Marshal
	default:
		// default use json
		reader.marshal = json.Marshal
		c.Writer.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
	}
	http.ServeContent(c.Writer, c.Request, name, modtime, reader)
}

type objectReader struct {
	obj     interface{}
	marshal func(v interface{}) ([]byte, error)
	reader  *bytes.Reader
}

func (r *objectReader) getReader() (reader io.ReadSeeker, e error) {
	if r.reader == nil {
		var b []byte
		b, e = r.marshal(r.obj)
		if e != nil {
			return
		}
		r.reader = bytes.NewReader(b)
	}
	reader = r.reader
	return
}
func (r *objectReader) Read(p []byte) (int, error) {
	reader, e := r.getReader()
	if e != nil {
		return 0, e
	}
	return reader.Read(p)
}
func (r *objectReader) Seek(offset int64, whence int) (int64, error) {
	reader, e := r.getReader()
	if e != nil {
		return 0, e
	}
	return reader.Seek(offset, whence)
}
