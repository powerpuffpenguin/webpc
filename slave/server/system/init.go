package system

import (
	"context"

	grpc_system "github.com/powerpuffpenguin/webpc/protocol/forward/system"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_system.RegisterSystemServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_system.RegisterSystemHandler(context.Background(), gateway, cc)
}
