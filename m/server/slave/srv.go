package slave

import (
	"context"
	"net/http"
	"strings"

	db0 "github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	"github.com/powerpuffpenguin/webpc/m/server/slave/internal/db"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type server struct {
	grpc_slave.UnimplementedSlaveServer
	helper.Helper
}

var emptyFindResponse grpc_slave.FindResponse

func (s server) Find(ctx context.Context, req *grpc_slave.FindRequest) (resp *grpc_slave.FindResponse, e error) {
	TAG := `slave Find`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	s.SetHTTPCacheMaxAge(ctx, 60)
	e = s.ServeMessage(ctx, db.LastModified(), func(nobody bool) error {
		if nobody {
			resp = &emptyFindResponse
		} else {
			tmp, err := db.Find(ctx, req)
			if err != nil {
				return err
			}
			if !userdata.AuthAny(db0.Root) {
				for _, item := range tmp.Data {
					item.Code = ``
				}
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
func (s server) Add(ctx context.Context, req *grpc_slave.AddRequest) (resp *grpc_slave.AddResponse, e error) {
	TAG := `slave Add`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == `` {
		e = s.Error(codes.InvalidArgument, `invalid name`)
		return
	}
	req.Description = strings.TrimSpace(req.Description)

	id, code, e := db.Add(req.Name, req.Description)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.String(`name`, req.Name),
				zap.String(`description`, req.Description),
			)
		}
		return
	}

	resp = &grpc_slave.AddResponse{
		Id:   id,
		Code: code,
	}
	s.SetHTTPCode(ctx, http.StatusCreated)
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`who`, userdata.Who()),
			zap.Int64(`id`, id),
			zap.String(`name`, req.Name),
			zap.String(`description`, req.Description),
			zap.String(`code`, code),
		)
	}
	return
}

func (s server) Code(ctx context.Context, req *grpc_slave.CodeRequest) (resp *grpc_slave.CodeResponse, e error) {
	TAG := `slave Code`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	changed, code, e := db.Code(req.Id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
			)
		}
		return
	} else if changed {
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.String(`code`, code),
			)
		}
	}
	resp = &grpc_slave.CodeResponse{
		Changed: changed,
		Code:    code,
	}
	return
}

var (
	changedChangeResponse    = grpc_slave.ChangeResponse{Changed: true}
	notChangedChangeResponse grpc_slave.ChangeResponse
)

func (s server) Change(ctx context.Context, req *grpc_slave.ChangeRequest) (resp *grpc_slave.ChangeResponse, e error) {
	TAG := `slave Change`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == `` {
		e = s.Error(codes.InvalidArgument, `invalid name`)
		return
	}
	req.Description = strings.TrimSpace(req.Description)
	changed, e := db.Change(req.Id, req.Name, req.Description)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.String(`name`, req.Name),
				zap.String(`description`, req.Description),
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
				zap.String(`name`, req.Name),
				zap.String(`description`, req.Description),
			)
		}
	} else {
		resp = &notChangedChangeResponse
	}
	return
}
func (s server) Remove(ctx context.Context, req *grpc_slave.RemoveRequest) (resp *grpc_slave.RemoveResponse, e error) {
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
	resp = &grpc_slave.RemoveResponse{
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
