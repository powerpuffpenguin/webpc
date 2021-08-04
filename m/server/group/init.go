package group

import (
	"context"

	grpc_group "github.com/powerpuffpenguin/webpc/protocol/group"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_group.RegisterGroupServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_group.RegisterGroupHandler(context.Background(), gateway, cc)
}
