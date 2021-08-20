package shell

import (
	"sync"
)

var defaultManager = Manager{
	keys: make(map[string]*Element),
}

func DefaultManager() *Manager {
	return &defaultManager
}

type ListInfo struct {
	ID       int64
	Name     string
	Attached bool
}
type Manager struct {
	mutex sync.Mutex
	keys  map[string]*Element
}

func (m *Manager) Unattach(username string, shellid int64) (e error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	element, ok := m.keys[username]
	if !ok {
		return
	}
	element.Unattach(shellid)
	return
}
func (m *Manager) Attach(conn Conn, username string, shellid int64, cols, rows uint16, newshell bool) (s *Shell, e error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	element, ok := m.keys[username]
	if !ok {
		element = &Element{
			keys: make(map[int64]*Shell),
		}
	}

	s, e = element.Attach(conn, username, shellid, cols, rows, newshell)
	if e != nil {
		return
	}

	if !ok {
		m.keys[username] = element
	}
	return
}
func (m *Manager) List(username string) (arrs []ListInfo) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if element, ok := m.keys[username]; ok {
		for _, s := range element.keys {
			arrs = append(arrs, ListInfo{
				ID:       s.shellid,
				Name:     s.name,
				Attached: s.IsAttack(),
			})
		}
	}
	return
}
