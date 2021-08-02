package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/andybalholm/brotli"

	"github.com/gin-gonic/gin"
)

var (
	// DefaultExcludedExtentions .
	DefaultExcludedExtentions = NewExcludedExtensions([]string{
		".png", ".gif", ".jpeg", ".jpg",
	})
	// DefaultOptions .
	DefaultOptions = &Options{
		ExcludedExtensions: DefaultExcludedExtentions,
	}
)

// Options .
type Options struct {
	ExcludedExtensions   ExcludedExtensions
	ExcludedPaths        ExcludedPaths
	ExcludedPathesRegexs ExcludedPathesRegexs
	GzDecompressFn       func(c *gin.Context)
	BrDecompressFn       func(c *gin.Context)
}

// Option .
type Option func(*Options)

// WithExcludedExtensions .
func WithExcludedExtensions(args []string) Option {
	return func(o *Options) {
		o.ExcludedExtensions = NewExcludedExtensions(args)
	}
}

// WithExcludedPaths .
func WithExcludedPaths(args []string) Option {
	return func(o *Options) {
		o.ExcludedPaths = NewExcludedPaths(args)
	}
}

// WithExcludedPathsRegexs .
func WithExcludedPathsRegexs(args []string) Option {
	return func(o *Options) {
		o.ExcludedPathesRegexs = NewExcludedPathesRegexs(args)
	}
}

// WithGzDecompressFn .
func WithGzDecompressFn(decompressFn func(c *gin.Context)) Option {
	return func(o *Options) {
		o.GzDecompressFn = decompressFn
	}
}

// WithBrDecompressFn .
func WithBrDecompressFn(decompressFn func(c *gin.Context)) Option {
	return func(o *Options) {
		o.BrDecompressFn = decompressFn
	}
}

// ExcludedExtensions Using map for better lookup performance
type ExcludedExtensions map[string]bool

// NewExcludedExtensions .
func NewExcludedExtensions(extensions []string) ExcludedExtensions {
	res := make(ExcludedExtensions)
	for _, e := range extensions {
		res[e] = true
	}
	return res
}

// Contains .
func (e ExcludedExtensions) Contains(target string) bool {
	_, ok := e[target]
	return ok
}

// ExcludedPaths .
type ExcludedPaths []string

// NewExcludedPaths .
func NewExcludedPaths(paths []string) ExcludedPaths {
	return ExcludedPaths(paths)
}

// Contains .
func (e ExcludedPaths) Contains(requestURI string) bool {
	for _, path := range e {
		if strings.HasPrefix(requestURI, path) {
			return true
		}
	}
	return false
}

// ExcludedPathesRegexs .
type ExcludedPathesRegexs []*regexp.Regexp

// NewExcludedPathesRegexs .
func NewExcludedPathesRegexs(regexs []string) ExcludedPathesRegexs {
	result := make([]*regexp.Regexp, len(regexs), len(regexs))
	for i, reg := range regexs {
		result[i] = regexp.MustCompile(reg)
	}
	return result
}

// Contains .
func (e ExcludedPathesRegexs) Contains(requestURI string) bool {
	for _, reg := range e {
		if reg.MatchString(requestURI) {
			return true
		}
	}
	return false
}

// GzDefaultDecompressHandle .
func GzDefaultDecompressHandle(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}
	r, err := gzip.NewReader(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Request.Header.Del("Content-Encoding")
	c.Request.Header.Del("Content-Length")
	c.Request.Body = r
}

// BrDefaultDecompressHandle .
func BrDefaultDecompressHandle(c *gin.Context) {
	if c.Request.Body == nil {
		return
	}

	r := _ReaderCloser{
		brotli.NewReader(c.Request.Body),
		c.Request.Body,
	}
	c.Request.Header.Del("Content-Encoding")
	c.Request.Header.Del("Content-Length")
	c.Request.Body = r
}

type _ReaderCloser struct {
	io.Reader
	io.Closer
}
