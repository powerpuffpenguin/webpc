package configure

import (
	"path/filepath"

	"github.com/powerpuffpenguin/webpc/utils"
)

type System struct {
	Enable bool
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
