package group

import (
	"context"
	"net/http"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	"github.com/powerpuffpenguin/webpc/signal/group"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"

	grpc_group "github.com/powerpuffpenguin/webpc/protocol/group"
)

type server struct {
	grpc_group.UnimplementedGroupServer
	helper.Helper
}

var emptyListResponse grpc_group.ListResponse

func (s server) List(ctx context.Context, req *grpc_group.ListRequest) (resp *grpc_group.ListResponse, e error) {
	TAG := `groud List`
	t := group.DefaultTree()
	s.SetHTTPCacheMaxAge(ctx, 5)
	e = s.ServeMessage(ctx, t.LastModified(), func(nobody bool) error {
		if nobody {
			resp = &emptyListResponse
		} else {
			items := make([]*grpc_group.Data, 0, t.Len())
			err := t.Foreach(func(ele *group.Element) (e error) {
				var children []int64
				if len(ele.Children) != 0 {
					children = make([]int64, len(ele.Children))
					for i, child := range ele.Children {
						children[i] = child.ID
					}
				}
				items = append(items, &grpc_group.Data{
					Id:          ele.ID,
					Name:        ele.Name,
					Description: ele.Description,
					Children:    children,
				})
				return nil
			})
			if err != nil {
				return nil
			}
			resp = &grpc_group.ListResponse{
				Items: items,
			}
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
func (s server) Add(ctx context.Context, req *grpc_group.AddRequest) (resp *grpc_group.AddResponse, e error) {
	if req.Parent == 0 {
		e = s.Error(codes.InvalidArgument, `parent not supported: 0`)
		return
	} else if req.Name == `` {
		e = s.Error(codes.InvalidArgument, `name not supported empty`)
		return
	}

	TAG := `group Add`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	id, e := group.DefaultTree().Add(ctx, req.Parent, req.Name, req.Description)
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

	resp = &grpc_group.AddResponse{
		Id: id,
	}
	s.SetHTTPCode(ctx, http.StatusCreated)
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`who`, userdata.Who()),
			zap.Int64(`id`, id),
			zap.String(`name`, req.Name),
			zap.String(`description`, req.Description),
		)
	}
	return
}

var (
	changedMoveResponse = grpc_group.MoveResponse{
		Changed: true,
	}
	notChangedMoveResponse = grpc_group.MoveResponse{}
)

func (s server) Move(ctx context.Context, req *grpc_group.MoveRequest) (resp *grpc_group.MoveResponse, e error) {
	TAG := `group Move`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	changed, e := group.DefaultTree().Move(ctx, req.Id, req.Parent)
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
	} else if changed {
		resp = &changedMoveResponse
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.Int64(`parent`, req.Parent),
			)
		}
	} else {
		resp = &notChangedMoveResponse
	}
	return
}

var (
	changedChangeResponse    = grpc_group.ChangeResponse{Changed: true}
	notChangedChangeResponse = grpc_group.ChangeResponse{}
)

func (s server) Change(ctx context.Context, req *grpc_group.ChangeRequest) (resp *grpc_group.ChangeResponse, e error) {
	TAG := `group Change`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	changed, e := group.DefaultTree().Change(ctx, req.Id, req.Name, req.Description)
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
	} else if changed {
		resp = &changedChangeResponse
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
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

var emptyRemoveResponse grpc_group.RemoveResponse

func (s server) Remove(ctx context.Context, req *grpc_group.RemoveRequest) (resp *grpc_group.RemoveResponse, e error) {
	TAG := `group Remove`
	_, userdata, e := s.Userdata(ctx)
	if e != nil {
		return
	}

	rowsAffected, e := group.DefaultTree().Remove(ctx, req.Id)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
			)
		}
		return
	} else if rowsAffected == 0 {
		resp = &emptyRemoveResponse
	} else {
		resp = &grpc_group.RemoveResponse{
			RowsAffected: int64(rowsAffected),
		}
		if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`who`, userdata.Who()),
				zap.Int64(`id`, req.Id),
				zap.Int(`rowsAffected`, rowsAffected),
			)
		}
	}
	return
}
