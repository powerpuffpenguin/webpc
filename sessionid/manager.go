package sessionid

import (
	"context"
	"strconv"

	"github.com/powerpuffpenguin/sessionstore"
	"github.com/powerpuffpenguin/sessionstore/cryptoer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Manager struct {
	m         *sessionstore.Manager
	platforms map[string]bool
}

func (m *Manager) getError(e error) error {
	if e == cryptoer.ErrExpired || e == cryptoer.ErrNotExistsToken {
		return status.Error(codes.Unauthenticated, e.Error())
	} else if e == cryptoer.ErrHashUnavailable ||
		e == cryptoer.ErrInvalidToken ||
		e == cryptoer.ErrRefreshTokenNotMatched ||
		e == cryptoer.ErrCannotRefresh {
		return status.Error(codes.PermissionDenied, e.Error())
	}
	return e
}
func (m *Manager) Get(ctx context.Context, access string) (session *Session, e error) {
	token, s, e := m.m.Get(ctx, access)
	if e != nil {
		e = m.getError(e)
		return
	}
	session = s.(*Session)
	session.Token = token
	return
}
func (m *Manager) Put(ctx context.Context, id int64, platform string, session *Session) (token *sessionstore.Token, e error) {
	if !m.platforms[platform] {
		e = status.Error(codes.InvalidArgument, `platform not supported`)
		return
	}

	token, e = m.m.Put(ctx, strconv.FormatInt(id, 10), platform, session)
	if e != nil {
		e = m.getError(e)
		return
	}
	session.Token = token
	return
}
func (m *Manager) Delete(ctx context.Context, access string) (e error) {
	e = m.m.Delete(ctx, access)
	if e != nil {
		e = m.getError(e)
		return
	}
	return
}
func (m *Manager) DeleteID(ctx context.Context, id int64) (e error) {
	e = m.m.DeleteID(ctx, strconv.FormatInt(id, 10))
	if e != nil {
		e = m.getError(e)
		return
	}
	return
}
func (m *Manager) Refresh(ctx context.Context, access, refresh string) (token *sessionstore.Token, e error) {
	token, _, e = m.m.Refresh(ctx, access, refresh)
	if e != nil {
		e = m.getError(e)
		return
	}
	return
}
