package fs

import (
	"context"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_fs.RegisterFSServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_fs.RegisterFSHandler(context.Background(), gateway, cc)
}
