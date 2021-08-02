package helper

import (
	"context"
	"errors"
	"strings"
	"github.com/powerpuffpenguin/webpc/sessions"

	"github.com/powerpuffpenguin/sessionid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	if access != `` {
		session, e = sessions.DefaultManager().Get(ctx, access)
		if e != nil {
			e = h.ToTokenError(e)
		}
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
func (h Helper) Userdata(ctx context.Context, prepare ...string) (newctx context.Context, userdata sessions.Userdata, e error) {
	newctx, session, e := h.Session(ctx)
	if e != nil {
		return
	}
	if len(prepare) != 0 {
		e = session.Prepare(ctx, prepare...)
		if e != nil {
			e = h.ToTokenError(e)
			return
		}
	}
	e = session.Get(newctx, sessions.KeyUserdata, &userdata)
	if e != nil {
		e = h.ToTokenError(e)
	}
	return
}

func (h Helper) ToTokenError(e error) error {
	if sessionid.IsTokenExpired(e) {
		e = h.Error(codes.Unauthenticated, e.Error())
	} else if errors.Is(e, sessionid.ErrTokenNotExists) {
		e = h.Error(codes.PermissionDenied, `token not exists`)
	}
	return e
}
