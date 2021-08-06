package session

import (
	"context"
	"net/http"
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_session "github.com/powerpuffpenguin/webpc/protocol/session"
	"github.com/powerpuffpenguin/webpc/sessions"
	signal_session "github.com/powerpuffpenguin/webpc/signal/session"

	"github.com/powerpuffpenguin/sessionid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type server struct {
	grpc_session.UnimplementedSessionServer
	helper.Helper
}

func (s server) Signin(ctx context.Context, req *grpc_session.SigninRequest) (resp *grpc_session.SigninResponse, e error) {
	TAG := `session Signin`
	at := time.Unix(req.Unix, 0)
	now := time.Now()
	if at.After(now.Add(time.Minute*5)) ||
		at.Before(now.Add(-time.Minute*5)) {
		e = s.Error(codes.InvalidArgument, `The request has expired, please confirm that the system time is correct.`)
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`name`, req.Name),
				zap.String(`password`, req.Password),
				zap.Int64(`unix`, req.Unix),
				zap.String(`at`, at.Local().String()),
			)
		}
		return
	}
	result, e := signal_session.Signin(ctx, req.Platform, req.Unix, req.Name, req.Password)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`name`, req.Name),
				zap.String(`password`, req.Password),
				zap.Int64(`unix`, req.Unix),
				zap.String(`at`, at.Local().String()),
			)
		}
		return
	} else if result.ID == 0 {
		e = s.Error(codes.NotFound, `name or password not match`)
		return
	}
	id, e := sessions.PlatformID(req.Platform, result.ID)
	if e != nil {
		e = s.Error(codes.InvalidArgument, e.Error())
		return
	}
	m := sessions.DefaultManager()
	e = m.Destroy(ctx, id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`name`, req.Name),
				zap.String(`password`, req.Password),
				zap.Int64(`unix`, req.Unix),
				zap.String(`at`, at.Local().String()),
			)
		}
		return
	}
	session, refresh, e := m.Create(ctx,
		id,
		sessionid.Pair{
			Key: sessions.KeyUserdata,
			Value: &sessions.Userdata{
				ID:            result.ID,
				Parent:        result.Parent,
				Name:          result.Name,
				Nickname:      result.Nickname,
				Authorization: result.Authorization,
			},
		},
		sessionid.Pair{
			Key:   sessions.KeyModtime,
			Value: time.Now(),
		},
	)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`name`, req.Name),
				zap.String(`password`, req.Password),
				zap.Int64(`unix`, req.Unix),
				zap.String(`at`, at.Local().String()),
			)
		}
		return
	}

	access := session.Token()
	if req.Cookie {
		s.SetHTTPCookie(ctx, &http.Cookie{
			Name:     s.CookieName(),
			Value:    access,
			Path:     `/`,
			HttpOnly: true,
		})
	}

	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.Error(e),
			zap.String(`access`, access),
			zap.String(`refresh`, refresh),
			zap.Int64(`id`, result.ID),
			zap.String(`name`, req.Name),
			zap.String(`password`, req.Password),
			zap.Int64(`unix`, req.Unix),
			zap.String(`at`, at.Local().String()),
		)
	}
	resp = &grpc_session.SigninResponse{
		Access:  access,
		Refresh: refresh,

		Id:            result.ID,
		Name:          result.Name,
		Nickname:      result.Nickname,
		Authorization: result.Authorization,
	}
	return
}

var emptySignoutResponse grpc_session.SignoutResponse

func (s server) Signout(ctx context.Context, req *grpc_session.SignoutRequest) (resp *grpc_session.SignoutResponse, e error) {
	TAG := `session Signout`
	access := s.GetToken(ctx)
	if access == `` {
		resp = &emptySignoutResponse
		return
	}
	e = sessions.DefaultManager().DestroyByToken(ctx, access)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, access),
			)
		}
		return
	}
	if req.Cookie {
		s.SetHTTPCookie(ctx, &http.Cookie{
			Name:     s.CookieName(),
			Value:    ``,
			Path:     `/`,
			HttpOnly: true,
			MaxAge:   -1,
		})
	}
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`access`, access),
		)
	}
	resp = &emptySignoutResponse
	return
}
func (s server) Refresh(ctx context.Context, req *grpc_session.RefreshRequest) (resp *grpc_session.RefreshResponse, e error) {
	TAG := `session Refresh`
	access, refresh, e := sessions.DefaultManager().Refresh(ctx, req.Access, req.Refresh)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`access`, req.Access),
				zap.String(`refresh`, req.Refresh),
			)
		}
		return
	}
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`access`, req.Access),
			zap.String(`refresh`, req.Refresh),
			zap.String(`new access`, access),
			zap.String(`new refresh`, refresh),
		)
	}
	resp = &grpc_session.RefreshResponse{
		Access:  access,
		Refresh: refresh,
	}
	return
}

var truePasswordResponse = grpc_session.PasswordResponse{
	Changed: true,
}
var falsePasswordResponse grpc_session.PasswordResponse

func (s server) Password(ctx context.Context, req *grpc_session.PasswordRequest) (resp *grpc_session.PasswordResponse, e error) {
	TAG := `session Password`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	changed, e := signal_session.Password(ctx, userdata.ID, req.Old, req.Password)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.Int64(`id`, userdata.ID),
				zap.String(`name`, userdata.Name),
			)
		}
		return
	}
	if changed {
		resp = &truePasswordResponse
		ed := sessions.DestroyID(userdata.ID)
		if ed != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
				ce.Write(
					zap.Error(ed),
					zap.Int64(`id`, userdata.ID),
					zap.String(`name`, userdata.Name),
				)
			}
		}
	} else {
		resp = &falsePasswordResponse
	}
	return
}

var emptyUserResponse grpc_session.UserResponse

func (s server) User(ctx context.Context, req *grpc_session.UserRequest) (resp *grpc_session.UserResponse, e error) {
	TAG := `session User`
	_, session, e := s.Session(ctx)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	e = session.Prepare(ctx, sessions.KeyModtime, sessions.KeyUserdata)
	if e != nil {
		e = s.ToTokenError(e)
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	var modtime time.Time
	e = session.Get(ctx, sessions.KeyModtime, &modtime)
	if e != nil {
		e = s.ToTokenError(e)
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	s.SetHTTPCacheMaxAge(ctx, 60)
	e = s.ServeMessage(ctx, modtime, func(nobody bool) error {
		if nobody {
			resp = &emptyUserResponse
		} else {
			var userdata sessions.Userdata
			if e := session.Get(ctx, sessions.KeyUserdata, &userdata); e != nil {
				e = s.ToTokenError(e)
				if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
					ce.Write(
						zap.Error(e),
					)
				}
				return e
			}
			resp = &grpc_session.UserResponse{
				Id:            userdata.ID,
				Parent:        userdata.Parent,
				Name:          userdata.Name,
				Nickname:      userdata.Nickname,
				Authorization: userdata.Authorization,
			}
		}
		return nil
	})
	return
}
