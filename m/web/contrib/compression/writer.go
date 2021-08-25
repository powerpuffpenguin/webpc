package compression

import (
	"compress/gzip"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/gin-gonic/gin"
)

const (
	// GzBestCompression .
	GzBestCompression = gzip.BestCompression
	// GzBestSpeed .
	GzBestSpeed = gzip.BestSpeed
	// GzDefaultCompression .
	GzDefaultCompression = gzip.DefaultCompression
	// GzNoCompression .
	GzNoCompression = gzip.NoCompression

	// BrBestCompression .
	BrBestCompression = brotli.BestCompression
	// BrBestSpeed .
	BrBestSpeed = brotli.BestSpeed
	// BrDefaultCompression .
	BrDefaultCompression = brotli.DefaultCompression
)

// Compression .
func Compression(br, gz int, options ...Option) gin.HandlerFunc {
	return newHandler(br, gz, options...).Handle
}

type _Writer struct {
	gin.ResponseWriter
	w io.Writer
}

func (g *_Writer) WriteString(s string) (int, error) {
	return g.w.Write([]byte(s))
}

func (g *_Writer) Write(data []byte) (int, error) {
	n, e := g.w.Write(data)
	return n, e
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *_Writer) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}
