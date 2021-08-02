package db

import (
	"github.com/powerpuffpenguin/webpc/logger"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

var defaultFilesystem Filesystem

// DefaultFilesystem .
func DefaultFilesystem() *Filesystem {
	return &defaultFilesystem
}

// Filesystem .
type Filesystem struct {
	root, name, ext string
}

func (f *Filesystem) onStart(filename string) (e error) {
	root, filename := filepath.Split(filename)
	ext := filepath.Ext(filename)
	name := filename
	if len(ext) != 0 {
		name = name[:len(name)-len(ext)]
	}
	f.root = root
	f.name = name
	f.ext = ext
	if ce := logger.Logger.Check(zap.InfoLevel, `logger set fileroot`); ce != nil {
		ce.Write(
			zap.String(`fileroot`, root),
			zap.String(`name`, name),
			zap.String(`ext`, ext),
		)
	}
	return
}

// Get os sys filename
func (f *Filesystem) Get(id string) (filename string, allowed bool) {
	name := filepath.Clean(filepath.Join(f.root, id))
	if !strings.HasPrefix(name, f.root) ||
		!strings.HasPrefix(id, f.name) || !strings.HasSuffix(id, f.ext) {
		return
	}
	allowed = true
	filename = name
	return
}

// Stat return root directory stat
func (f *Filesystem) Stat() (stat os.FileInfo, e error) {
	stat, e = os.Stat(f.root)
	if e != nil {
		return
	}
	return
}

// List return all logs filename
func (f *Filesystem) List() (names []string, e error) {
	file, e := os.Open(f.root)
	if e != nil {
		return
	}
	stats, e := file.Readdir(0)
	file.Close()
	if e != nil {
		return
	}
	for _, stat := range stats {
		if stat.IsDir() {
			continue
		}
		name := stat.Name()
		if !strings.HasPrefix(name, f.name) || !strings.HasSuffix(name, f.ext) {
			continue
		}
		names = append(names, name)
	}
	return
}
