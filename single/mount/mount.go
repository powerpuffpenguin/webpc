package mount

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Separator = string(filepath.Separator)

type FileInfo struct {
	Name  string
	Mode  uint32
	Size  int64
	IsDir bool
}

type Mount struct {
	name                string
	root                string
	read, write, shared bool
}

func (m *Mount) Name() string {
	return m.name
}

func (m *Mount) Read() bool {
	return m.read
}

func (m *Mount) Write() bool {
	return m.write
}

func (m *Mount) Shared() bool {
	return m.shared
}
func (m *Mount) toError(name string, e error) error {
	if os.IsNotExist(e) {
		e = status.Error(codes.NotFound, `not exists: `+name)
	} else if os.IsExist(e) {
		e = status.Error(codes.AlreadyExists, `already exists: `+name)
	} else if os.IsPermission(e) {
		e = status.Error(codes.PermissionDenied, `forbidden: `+name)
	}
	return e
}
func (m *Mount) LS(path string) (dir string, modtime time.Time, results []FileInfo, e error) {
	dst, e := m.filename(path)
	if e != nil {
		return
	}
	f, e := os.Open(dst)
	if e != nil {
		e = m.toError(path, e)
		return
	}
	defer f.Close()
	stat, e := f.Stat()
	if e != nil {
		e = m.toError(path, e)
		return
	}
	modtime = stat.ModTime()

	infos, e := f.Readdir(0)
	count := len(infos)
	if e != nil {
		if count == 0 {
			return
		}
		if ce := logger.Logger.Check(zap.WarnLevel, "readdir error"); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		e = nil
	}

	dir = dst[len(m.root):]
	if Separator != `/` {
		dir = strings.ReplaceAll(path, Separator, `/`)
	}
	if !strings.HasPrefix(dir, `/`) {
		dir = `/` + dir
	}
	results = make([]FileInfo, count)
	for i := 0; i < count; i++ {
		results[i].Name = infos[i].Name()
		results[i].Mode = uint32(infos[i].Mode())
		results[i].Size = infos[i].Size()
		results[i].IsDir = infos[i].IsDir()
	}
	return
}

func (m *Mount) filename(path string) (filename string, e error) {
	filename = filepath.Clean(m.root + path)
	if m.root != filename {
		root := m.root
		if !strings.HasSuffix(root, Separator) {
			root += Separator
		}
		if !strings.HasPrefix(filename, root) {
			e = status.Error(codes.PermissionDenied, `Illegal path`)
			return
		}
	}
	return
}
func (m *Mount) Open(name string) (*os.File, error) {
	return m.OpenFile(name, os.O_RDONLY, 0)
}
func (m *Mount) OpenFile(name string, flag int, perm os.FileMode) (f *os.File, e error) {
	path, e := m.filename(name)
	if e != nil {
		return
	}
	f, e = os.OpenFile(path, flag, perm)
	if e != nil {
		e = m.toError(name, e)
		return
	}
	return
}
func (m *Mount) checkName(name string) error {
	val := filepath.Base(name)
	if name != val {
		return status.Error(codes.InvalidArgument, `invalid name: `+name)
	} else if name == `` {
		return status.Error(codes.InvalidArgument, `name not supported empty`)
	}
	return nil
}
func (m *Mount) formatCreate(dir, name string) (string, error) {
	e := m.checkName(name)
	if e != nil {
		return ``, e
	}
	dir, e = m.filename(dir)
	if e != nil {
		return ``, e
	}
	return filepath.Join(dir, name), nil
}
func (m *Mount) Create(file bool, dir, name string, perm os.FileMode) (info os.FileInfo, e error) {
	path, e := m.formatCreate(dir, name)
	if e != nil {
		return
	}
	if file {
		var f *os.File
		f, e = os.OpenFile(path, os.O_CREATE|os.O_EXCL, perm)
		if e != nil {
			e = m.toError(filepath.Join(dir, name), e)
			return
		}
		info, e = f.Stat()
		f.Close()
		if e != nil {
			e = m.toError(filepath.Join(dir, name), e)
			return
		}
	} else {
		e = os.Mkdir(path, perm)
		if e != nil {
			e = m.toError(filepath.Join(dir, name), e)
			return
		}
		info, e = os.Stat(path)
		if e != nil {
			e = m.toError(filepath.Join(dir, name), e)
			return
		}
	}
	return
}
func (m *Mount) RemoveAll(dir string, names []string) (e error) {
	if len(names) == 0 {
		return
	}
	dir, e = m.filename(dir)
	if e != nil {
		return
	}
	for _, name := range names {
		e = m.checkName(name)
		if e != nil {
			return
		}
	}
	for _, name := range names {
		dst := filepath.Join(dir, name)
		e = os.RemoveAll(dst)
		if e != nil {
			e = m.toError(dst, e)
			return
		}
	}
	return
}
func (m *Mount) Rename(dir, old, current string) (e error) {
	dir, e = m.filename(dir)
	if e != nil {
		return
	}
	e = m.checkName(old)
	if e != nil {
		return
	}
	e = m.checkName(current)
	if e != nil {
		return
	}

	old = filepath.Join(dir, old)
	current = filepath.Join(dir, current)
	e = os.Rename(old, current)
	if e != nil {
		if os.IsNotExist(e) {
			e = status.Error(codes.NotFound, `not exists: `+old)
		} else if os.IsExist(e) {
			e = status.Error(codes.AlreadyExists, `already exists: `+current)
		} else if os.IsPermission(e) {
			e = status.Error(codes.PermissionDenied, `forbidden: `+current)
		}
		return
	}
	return
}
