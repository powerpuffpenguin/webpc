package session

import (
	"context"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/powerpuffpenguin/sessionstore"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_session "github.com/powerpuffpenguin/webpc/protocol/session"
	"github.com/powerpuffpenguin/webpc/sessionid"
	signal_session "github.com/powerpuffpenguin/webpc/signal/session"

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
	session := &sessionid.Session{
		ID:            result.ID,
		Name:          result.Name,
		Nickname:      result.Nickname,
		Authorization: result.Authorization,
		Parent:        result.Parent,
	}
	token, e := sessionid.DefaultManager().Put(ctx, session.ID, req.Platform, session)
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
	if req.Cookie {
		s.SetHTTPCookie(ctx, &http.Cookie{
			Name:     s.CookieName(),
			Value:    token.Access,
			Path:     `/`,
			HttpOnly: true,
		})
	}

	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.Error(e),
			zap.String(`access`, token.Access),
			zap.String(`refresh`, token.Refresh),
			zap.Int64(`id`, result.ID),
			zap.String(`name`, req.Name),
			zap.String(`password`, req.Password),
			zap.Int64(`unix`, req.Unix),
			zap.String(`at`, at.Local().String()),
		)
	}
	resp = &grpc_session.SigninResponse{
		Token: &grpc_session.Token{
			Access:          token.Access,
			Refresh:         token.Refresh,
			AccessDeadline:  token.AccessDeadline,
			RefreshDeadline: token.RefreshDeadline,
			Deadline:        token.Deadline,
		},
		Data: &grpc_session.Data{
			Id:            session.ID,
			Name:          session.Name,
			Nickname:      session.Nickname,
			Authorization: session.Authorization,
			Parent:        session.Parent,
		},
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
	e = sessionid.DefaultManager().Delete(ctx, access)
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
	token, e := sessionid.DefaultManager().Refresh(ctx, req.Access, req.Refresh)
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
			zap.String(`new access`, token.Access),
			zap.String(`new refresh`, token.Refresh),
		)
	}
	resp = &grpc_session.RefreshResponse{
		Token: &grpc_session.Token{
			Access:          token.Access,
			Refresh:         token.Access,
			AccessDeadline:  token.AccessDeadline,
			RefreshDeadline: token.RefreshDeadline,
			Deadline:        token.Deadline,
		},
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
		err := sessionid.DefaultManager().DeleteID(context.Background(), userdata.ID)
		if err != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
				ce.Write(
					zap.Error(err),
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

	modtime := time.Unix(session.Token.AccessDeadline, 0)
	e = s.ServeMessage(ctx, modtime, func(nobody bool) error {
		if nobody {
			resp = &emptyUserResponse
		} else {
			s.SetHTTPCacheExpress(ctx, modtime)
			resp = &grpc_session.UserResponse{
				Id:            session.ID,
				Parent:        session.Parent,
				Name:          session.Name,
				Nickname:      session.Nickname,
				Authorization: session.Authorization,
			}
		}
		return nil
	})
	return
}
func (s server) Download(ctx context.Context, req *grpc_session.DownloadRequest) (resp *grpc_session.DownloadResponse, e error) {
	_, session, _ := s.Session(ctx)
	if session == nil {
		session = &sessionid.Session{}
	}
	unix := time.Now().Add(time.Hour * 24).Unix()
	session.Token = &sessionstore.Token{
		Access:          `temporary`,
		Refresh:         `temporary`,
		AccessDeadline:  unix,
		RefreshDeadline: unix,
		Deadline:        unix,
	}
	b, e := session.Marshal()
	if e != nil {
		return
	}
	playdata := base64.RawURLEncoding.EncodeToString(b)
	access, e := sessionid.DefaultManager().Sin(playdata)
	if e != nil {
		return
	}
	resp = &grpc_session.DownloadResponse{
		Access: access,
	}
	return
}
