package user

import (
	"context"
	"github.com/powerpuffpenguin/webpc/m/server/user/internal/db"
	grpc_user "github.com/powerpuffpenguin/webpc/protocol/user"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_user.RegisterUserServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	db.Init()
	return grpc_user.RegisterUserHandler(context.Background(), gateway, cc)
}
