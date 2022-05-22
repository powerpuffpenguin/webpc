package helper

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/powerpuffpenguin/sessionstore/cryptoer"
	"github.com/powerpuffpenguin/webpc/sessionid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type sessionKey struct{}
type sessionValue struct {
	session *sessionid.Session
	e       error
}

func (h Helper) GetToken(ctx context.Context) (token string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		strs := md.Get(`Authorization`)
		for _, str := range strs {
			if strings.HasPrefix(str, `Bearer `) {
				token = str[7:]
				break
			}
		}
	}
	return
}
func (h Helper) accessSession(ctx context.Context) (session *sessionid.Session, e error) {
	access := h.GetToken(ctx)
	if access == `` {
		e = h.Error(codes.PermissionDenied, `not found token`)
	} else {
		session, e = sessionid.DefaultManager().Get(ctx, access)
	}
	return
}
func (h Helper) Session(ctx context.Context) (newctx context.Context, session *sessionid.Session, e error) {
	newctx = ctx

	cache, ok := ctx.Value(sessionKey{}).(sessionValue)
	if ok {
		session = cache.session
		e = cache.e
		return
	}
	session, e = h.accessSession(ctx)
	if e == nil && session == nil {
		e = h.Error(codes.PermissionDenied, `token not exists`)
	}
	newctx = context.WithValue(ctx, sessionKey{}, sessionValue{
		session: session,
		e:       e,
	})
	return
}
func (h Helper) Userdata(ctx context.Context) (newctx context.Context, session *sessionid.Session, e error) {
	return h.Session(ctx)
}

type internalKey struct{}

func (h Helper) internalAccessSession(ctx context.Context) (session *sessionid.Session, e error) {
	access := h.GetToken(ctx)
	if access == `` {
		e = h.Error(codes.PermissionDenied, codes.PermissionDenied.String())
	} else if access == `Expired` {
		e = h.Error(codes.Unauthenticated, cryptoer.ErrExpired.Error())
	} else {
		var b []byte
		b, e = base64.RawURLEncoding.DecodeString(access)
		if e != nil {
			e = h.Error(codes.PermissionDenied, e.Error())
			return
		}
		var s sessionid.Session
		e = s.Unmarshal(b)
		if e != nil {
			e = h.Error(codes.PermissionDenied, e.Error())
			return
		}
		session = &s
		if session.Token.IsDeleted() {
			e = status.Error(codes.Unauthenticated, cryptoer.ErrNotExistsToken.Error())
			return
		} else if session.Token.IsExpired() {
			e = status.Error(codes.Unauthenticated, cryptoer.ErrExpired.Error())
			return
		}
	}
	return
}
func (h Helper) InternalSession(ctx context.Context) (newctx context.Context, session *sessionid.Session, e error) {
	newctx = ctx

	cache, ok := ctx.Value(internalKey{}).(sessionValue)
	if ok {
		session = cache.session
		return
	}
	session, e = h.internalAccessSession(ctx)
	newctx = context.WithValue(ctx, internalKey{}, sessionValue{
		session: session,
		e:       e,
	})
	return
}
func (h Helper) JSONUserdata(ctx context.Context) (context.Context, *sessionid.Session, error) {
	return h.InternalSession(ctx)
}
