package shell

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/db"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	fmt.Println(`Unattach`)
	defer fmt.Println(`Unattach end`)
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
func (m *Manager) restore() {
	items, e := db.List()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `list shell from db error`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	for _, item := range items {
		m.restoreItem(&item)
	}
}
func (m *Manager) restoreItem(item *db.DataOfSlaveShell) {
	s := &Shell{
		term:       newTerm(item.UserName),
		username:   item.UserName,
		shellid:    item.ID,
		name:       item.Name,
		fontSize:   item.FontSize,
		fontFamily: item.FontFamily,
	}
	e := s.Run(nil, 24, 10)
	if e != nil {
		if ce := logger.Logger.Check(zap.ErrorLevel, `restore shell error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`username`, item.UserName),
				zap.Int64(`id`, item.ID),
				zap.String(`name`, item.Name),
			)
		}
		return
	}

	if ce := logger.Logger.Check(zap.InfoLevel, `restore shell`); ce != nil {
		ce.Write(
			zap.String(`username`, item.UserName),
			zap.Int64(`id`, item.ID),
			zap.String(`name`, item.Name),
		)
	}

	element := m.keys[item.UserName]
	if element == nil {
		element = &Element{
			keys: make(map[int64]*Shell),
		}
		m.keys[item.UserName] = element
	}
	element.keys[s.shellid] = s
}
func (m *Manager) Rename(username string, shellid int64, name string) (e error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	element, ok := m.keys[username]
	if !ok {
		e = status.Error(codes.NotFound, `shellid not exists: `+strconv.FormatInt(shellid, 10))
		return
	}

	e = element.Rename(shellid, name)
	return
}
func (m *Manager) Kill(username string, shellid int64) (e error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	element, ok := m.keys[username]
	if !ok {
		e = status.Error(codes.NotFound, `shellid not exists: `+strconv.FormatInt(shellid, 10))
		return
	}

	e = element.Kill(shellid)
	return
}
