package sessionid

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/powerpuffpenguin/sessionstore"
	grpc_session "github.com/powerpuffpenguin/webpc/protocol/session"
)

const (
	KeyUserdata = `userdata`
	KeyModtime  = `modtime`
)

type Session struct {
	ID            int64
	Name          string
	Nickname      string
	Authorization []int32
	Parent        int64

	Token *sessionstore.Token
}

func (d *Session) Who() string {
	return fmt.Sprintf("id=%d name=%s nickname=%s", d.ID, d.Name, d.Nickname)
}

func (d *Session) Unmarshal(b []byte) (e error) {
	var m grpc_session.Raw
	e = proto.Unmarshal(b, &m)
	if e != nil {
		return
	}
	s := m.Data
	if s == nil {
		e = errors.New(`token raw.data nil`)
		return
	}
	t := m.Token
	if t == nil {
		e = errors.New(`token raw.token nil`)
		return
	}
	d.ID = s.Id
	d.Name = s.Name
	d.Nickname = s.Nickname
	d.Authorization = s.Authorization
	d.Parent = s.Parent

	d.Token = sessionstore.NewToken(
		t.Access, t.Refresh,
		t.AccessDeadline, t.RefreshDeadline,
		t.Deadline,
	)
	return
}
func (d *Session) Marshal() (b []byte, e error) {
	if d.Token == nil {
		e = errors.New(`token nil`)
		return
	}
	m := &grpc_session.Raw{
		Data: &grpc_session.Data{
			Id:            d.ID,
			Name:          d.Name,
			Nickname:      d.Nickname,
			Authorization: d.Authorization,
			Parent:        d.Parent,
		},
		Token: &grpc_session.Token{
			Access:          d.Token.Access,
			Refresh:         d.Token.Refresh,
			AccessDeadline:  d.Token.AccessDeadline,
			RefreshDeadline: d.Token.RefreshDeadline,
			Deadline:        d.Token.Deadline,
		},
	}
	b, e = proto.Marshal(m)
	return
}

// Test if has all authorization return true
func (d *Session) AuthTest(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	var found bool
	for _, val := range vals {
		found = false
		for _, auth := range d.Authorization {
			if val == auth {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// AuthAny if has any authorization return true
func (d *Session) AuthAny(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	for _, val := range vals {
		for _, auth := range d.Authorization {
			if val == auth {
				return true
			}
		}
	}
	return false
}

// AuthNone if not has any authorization return true
func (d *Session) AuthNone(vals ...int32) bool {
	if d.ID == 0 {
		return false
	}
	for _, val := range vals {
		for _, auth := range d.Authorization {
			if val == auth {
				return false
			}
		}
	}
	return true
}

type Coder struct {
}

func (Coder) Unmarshal(b []byte) (session interface{}, e error) {
	var m grpc_session.Data
	e = proto.Unmarshal(b, &m)
	if e != nil {
		return
	}
	session = &Session{
		ID:            m.Id,
		Name:          m.Name,
		Nickname:      m.Nickname,
		Authorization: m.Authorization,
		Parent:        m.Parent,
	}
	return
}
func (Coder) Marshal(session interface{}) (b []byte, e error) {
	userdata := session.(*Session)
	m := &grpc_session.Data{
		Id:            userdata.ID,
		Name:          userdata.Name,
		Nickname:      userdata.Nickname,
		Authorization: userdata.Authorization,
		Parent:        userdata.Parent,
	}
	b, e = proto.Marshal(m)
	return
}
