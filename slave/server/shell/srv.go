package shell

import (
	"context"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/connect"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/shell"
	"go.uber.org/zap"
)

type server struct {
	grpc_shell.UnimplementedShellServer
	helper.Helper
}

func (s server) List(ctx context.Context, req *grpc_shell.ListRequest) (resp *grpc_shell.ListResponse, e error) {
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	items := shell.DefaultManager().List(userdata.Name)

	resp = &grpc_shell.ListResponse{
		Result: make([]*grpc_shell.ListResult, len(items)),
	}
	for i, item := range items {
		resp.Result[i] = &grpc_shell.ListResult{
			Id:       item.ID,
			Name:     item.Name,
			Attached: item.Attached,
		}
	}
	return
}

var emptyRenameResponse grpc_shell.RenameResponse

func (s server) Rename(ctx context.Context, req *grpc_shell.RenameRequest) (resp *grpc_shell.RenameResponse, e error) {
	TAG := `forward.shell Rename`
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	e = shell.DefaultManager().Rename(userdata.Name, req.Id, req.Name)
	if e != nil {
		return
	}
	resp = &emptyRenameResponse
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`who`, userdata.Who()),
			zap.Int64(`id`, req.Id),
			zap.String(`name`, req.Name),
		)
	}
	return
}

var emptyRemoveResponse grpc_shell.RemoveResponse

func (s server) Remove(ctx context.Context, req *grpc_shell.RemoveRequest) (resp *grpc_shell.RemoveResponse, e error) {
	TAG := `forward.shell Remove`
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	e = shell.DefaultManager().Kill(userdata.Name, req.Id)
	if e != nil {
		return
	}
	resp = &emptyRemoveResponse
	if ce := logger.Logger.Check(zap.InfoLevel, TAG); ce != nil {
		ce.Write(
			zap.String(`who`, userdata.Who()),
			zap.Int64(`id`, req.Id),
		)
	}
	return
}
func (s server) Connect(server grpc_shell.Shell_ConnectServer) (e error) {
	ctx := server.Context()
	_, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return
	}
	w := connect.New(server, userdata.Name)
	e = w.Serve()
	return
}
