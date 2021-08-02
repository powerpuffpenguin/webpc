package configure

import (
	"strings"
	"time"

	"github.com/powerpuffpenguin/webpc/utils"
)

type Session struct {
	Backend string
	Coder   string
	Client  SessionClient
	Memory  SessionMemory
}

func (s *Session) format() (e error) {
	s.Backend = strings.ToLower(strings.TrimSpace(s.Backend))
	if s.Backend == "memory" {
		e = s.Memory.format()
		if e != nil {
			return
		}
	}
	s.Coder = strings.ToUpper(strings.TrimSpace(s.Coder))
	return
}

type SessionClient struct {
	Protocol string
	Addr     string
	Token    string
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
	Access  time.Duration
	Refresh time.Duration
	MaxSize int
	Batch   int
	Clear   time.Duration
}
type SessionProviderRedis struct {
	URL     string
	Access  time.Duration
	Refresh time.Duration

	Batch       int
	KeyPrefix   string
	MetadataKey string
}
type SessionProviderBolt struct {
	Filename string

	Access  time.Duration
	Refresh time.Duration
	MaxSize int
	Batch   int
	Clear   time.Duration
}

func (s *SessionProviderBolt) format() (e error) {
	s.Filename = utils.Abs(utils.BasePath(), s.Filename)
	return
}
