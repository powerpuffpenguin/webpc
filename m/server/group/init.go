package group

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	tree "github.com/powerpuffpenguin/webpc/m/server/group/internal/tree"
	grpc_group "github.com/powerpuffpenguin/webpc/protocol/group"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_group.RegisterGroupServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	tree.Default().Init()
	return grpc_group.RegisterGroupHandler(context.Background(), gateway, cc)
}
