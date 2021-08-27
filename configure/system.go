package configure

import (
	"path/filepath"
	"runtime"

	"github.com/powerpuffpenguin/webpc/utils"
)

var defaultSystem *System

func DefaultSystem() *System {
	return defaultSystem
}

type System struct {
	Enable bool
	Shell  string
	VNC    string
	Mount  []Mount
}

func (s *System) format() (e error) {
	if !s.Enable {
		return
	}
	basePath := utils.BasePath()
	for i := 0; i < len(s.Mount); i++ {
		e = s.Mount[i].format(basePath)
		if e != nil {
			return
		}
	}
	if s.Shell == `` {
		if runtime.GOOS == "windows" {
			s.Shell = `shell-windows.bat`
		} else {
			s.Shell = `shell-linux`
		}
	}
	if filepath.IsAbs(s.Shell) {
		s.Shell = filepath.Clean(s.Shell)
	} else {
		s.Shell = filepath.Clean(filepath.Join(basePath, s.Shell))
	}
	return
}

type Mount struct {
	// web display name
	Name string
	// local file path
	Root string

	Read   bool
	Write  bool
	Shared bool
}

func (m *Mount) format(basePath string) (e error) {
	if filepath.IsAbs(m.Root) {
		m.Root = filepath.Clean(m.Root)
	} else {
		m.Root = filepath.Clean(filepath.Join(basePath, m.Root))
	}

	if m.Write || m.Shared {
		m.Read = true
	}
	return
}
