package user

import (
	"context"
	"net/http"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	"github.com/powerpuffpenguin/webpc/sessions"
	"github.com/powerpuffpenguin/webpc/utils"

	"github.com/powerpuffpenguin/webpc/m/server/user/internal/db"
	grpc_user "github.com/powerpuffpenguin/webpc/protocol/user"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type server struct {
	grpc_user.UnimplementedUserServer
	helper.Helper
}

var emptyFindResponse grpc_user.FindResponse

func (s server) Find(ctx context.Context, req *grpc_user.FindRequest) (resp *grpc_user.FindResponse, e error) {
	TAG := `user Find`
	s.SetHTTPCacheMaxAge(ctx, 5)
	e = s.ServeMessage(ctx, db.LastModified(), func(nobody bool) error {
		if nobody {
			resp = &emptyFindResponse
		} else {
			tmp, err := db.Find(ctx, req)
			if err != nil {
				return err
			}
			resp = tmp
		}
		return nil
	})
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
	}
	return
}
func (s server) Add(ctx context.Context, req *grpc_user.AddRequest) (resp *grpc_user.AddResponse, e error) {
	TAG := `user Add`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}
	if !utils.MatchName(req.Name) {
		e = s.Error(codes.InvalidArgument, `invalid name`)
		return
	} else if !utils.MatchPassword(req.Password) {
		e = s.Error(codes.InvalidArgument, `invalid password`)
		return
	}
	id, e := db.Add(ctx, req.Parent, req.Name, req.Nickname, req.Password, req.Authorization)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`parent`, req.Parent),
				zap.String(`name`, req.Name),
				zap.String(`nickname`, req.Nickname),
				zap.Int32s(`authorization`, req.Authorization),
			)
		}
		return
	}
	resp = &grpc_user.AddResponse{
		Id: id,
	}
	s.SetHTTPCode(ctx, http.StatusCreated)
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`who`, userdata.Who()),
			zap.Int64(`id`, id),
			zap.Int64(`parent`, req.Parent),
			zap.String(`name`, req.Name),
			zap.String(`nickname`, req.Nickname),
			zap.Int32s(`authorization`, req.Authorization),
		)
	}
	return
}

var (
	changedPasswordResponse    = grpc_user.PasswordResponse{Changed: true}
	notChangedPasswordResponse grpc_user.PasswordResponse
)

func (s server) Password(ctx context.Context, req *grpc_user.PasswordRequest) (resp *grpc_user.PasswordResponse, e error) {
	TAG := `user Password`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}
	if !utils.MatchPassword(req.Value) {
		e = s.Error(codes.InvalidArgument, `invalid password`)
		return
	}
	changed, e := db.Password(ctx, req.Id, req.Value)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
			)
		}
		return
	}
	if changed {
		resp = &changedPasswordResponse
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
			)
		}

		ed := sessions.DestroyID(req.Id)
		if ed != nil {
			if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
				ce.Write(
					zap.Error(ed),
					zap.String(`who`, userdata.Who()),
					zap.Int64(`id`, req.Id),
				)
			}
		}
	} else {
		resp = &notChangedPasswordResponse
	}
	return
}

var (
	changedChangeResponse    = grpc_user.ChangeResponse{Changed: true}
	notChangedChangeResponse grpc_user.ChangeResponse
)

func (s server) Change(ctx context.Context, req *grpc_user.ChangeRequest) (resp *grpc_user.ChangeResponse, e error) {
	TAG := `user Change`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}
	changed, e := db.Change(ctx, req.Id, req.Nickname, req.Authorization)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.String(`nickname`, req.Nickname),
				zap.Int32s(`authorization`, req.Authorization),
			)
		}
		return
	}
	if changed {
		resp = &changedChangeResponse
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.String(`nickname`, req.Nickname),
				zap.Int32s(`authorization`, req.Authorization),
			)
		}
	} else {
		resp = &notChangedChangeResponse
	}
	return
}

var (
	changedGroupResponse    = grpc_user.GroupResponse{Changed: true}
	notChangedGrouoResponse grpc_user.GroupResponse
)

func (s server) Group(ctx context.Context, req *grpc_user.GroupRequest) (resp *grpc_user.GroupResponse, e error) {
	TAG := `user Group`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}
	changed, e := db.Parent(ctx, req.Id, req.Parent)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.Int64(`parent`, req.Parent),
			)
		}
		return
	}
	if changed {
		resp = &changedGroupResponse
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.Int64(`parent`, req.Parent),
			)
		}
	} else {
		resp = &notChangedGrouoResponse
	}
	return
}
func (s server) Remove(ctx context.Context, req *grpc_user.RemoveRequest) (resp *grpc_user.RemoveResponse, e error) {
	TAG := `user Remove`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}
	rowsAffected, e := db.Remove(ctx, req.Id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64s(`id`, req.Id),
			)
		}
		return
	}
	resp = &grpc_user.RemoveResponse{
		RowsAffected: rowsAffected,
	}
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.Error(e),
			zap.String(`who`, userdata.Who()),
			zap.Int64s(`id`, req.Id),
			zap.Int64(`rowsAffected`, rowsAffected),
		)
	}
	return
}
