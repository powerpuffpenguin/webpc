package shell

import (
	"context"

	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_shell.RegisterShellServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_shell.RegisterShellHandler(context.Background(), gateway, cc)
}
