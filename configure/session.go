package configure

import (
	"strings"
	"time"

	"github.com/powerpuffpenguin/webpc/utils"
)

type Session struct {
	Memory SessionMemory
}

func (s *Session) format() (e error) {
	e = s.Memory.format()
	if e != nil {
		return
	}
	return
}

type SessionMemory struct {
	Manager  SessionManager
	Provider SessionProvider
}

func (s *SessionMemory) format() (e error) {
	s.Provider.Backend = strings.TrimSpace(strings.ToLower(s.Provider.Backend))
	if s.Provider.Backend == "bolt" {
		e = s.Provider.Bolt.format()
		if e != nil {
			return
		}
	}
	return
}

type SessionManager struct {
	Method string
	Key    string
}
type SessionProvider struct {
	Backend string
	Memory  SessionProviderMemory
	Redis   SessionProviderRedis
	Bolt    SessionProviderBolt
}
type SessionProviderMemory struct {
	Access   time.Duration
	Refresh  time.Duration
	Deadline time.Duration
	MaxSize  int
}
type SessionProviderRedis struct {
	URL      string
	Access   time.Duration
	Refresh  time.Duration
	Deadline time.Duration
}
type SessionProviderBolt struct {
	Filename string

	Access   time.Duration
	Refresh  time.Duration
	Deadline time.Duration
	MaxSize  int64
}

func (s *SessionProviderBolt) format() (e error) {
	s.Filename = utils.Abs(utils.BasePath(), s.Filename)
	return
}
