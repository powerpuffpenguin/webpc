package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

const (
	algorithmBrotli = iota + 1
	algorithmGzip
)

type helperHandler struct {
	*Options
	brPool *sync.Pool
	gzPool *sync.Pool
}

func newHandler(br, gz int, options ...Option) *helperHandler {
	var gzPool, brPool sync.Pool
	brPool.New = func() interface{} {
		br := brotli.NewWriterLevel(ioutil.Discard, br)
		return br
	}
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, gz)
		if err != nil {
			panic(err)
		}
		return gz
	}
	handler := &helperHandler{
		Options: DefaultOptions,
		brPool:  &brPool,
		gzPool:  &gzPool,
	}
	for _, setter := range options {
		setter(handler.Options)
	}
	return handler
}

func (h *helperHandler) Handle(c *gin.Context) {
	var encoding string
	if fn := h.BrDecompressFn; fn != nil {
		encoding = c.Request.Header.Get("Content-Encoding")
		if encoding == "br" {
			fn(c)
		}
	}
	if fn := h.GzDecompressFn; fn != nil {
		if encoding == "" {
			encoding = c.Request.Header.Get("Content-Encoding")
		}
		if encoding == "gzip" {
			fn(c)
		}
	}

	algorithm, yes := h.shouldCompress(c.Request)
	if !yes {
		return
	}
	var w io.WriteCloser
	switch algorithm {
	case algorithmGzip:
		gz := h.gzPool.Get().(*gzip.Writer)
		defer h.gzPool.Put(gz)
		defer gz.Reset(ioutil.Discard)
		gz.Reset(c.Writer)
		c.Header("Content-Encoding", "gzip")
		w = gz
	default:
		br := h.brPool.Get().(*brotli.Writer)
		defer h.brPool.Put(br)
		defer br.Reset(ioutil.Discard)
		br.Reset(c.Writer)
		c.Header("Content-Encoding", "br")
		w = br
	}
	c.Header("Vary", "Accept-Encoding")
	c.Writer = &_Writer{c.Writer, w}
	defer func() {
		w.Close()
		c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
	}()
	c.Next()
}

func (h *helperHandler) shouldCompress(req *http.Request) (algorithm int, yes bool) {
	encoding := req.Header.Get("Accept-Encoding")
	if strings.Contains(encoding, "br") {
		algorithm = algorithmBrotli
	} else if strings.Contains(encoding, "gzip") {
		algorithm = algorithmGzip
	} else {
		return
	}
	if strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/event-stream") {
		return
	}

	extension := filepath.Ext(req.URL.Path)
	if h.ExcludedExtensions.Contains(extension) {
		return
	}

	if h.ExcludedPaths.Contains(req.URL.Path) {
		return
	}
	if h.ExcludedPathesRegexs.Contains(req.URL.Path) {
		return
	}

	yes = true
	return
}
