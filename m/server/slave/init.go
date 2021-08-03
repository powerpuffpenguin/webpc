package slave

import (
	"context"

	"github.com/powerpuffpenguin/webpc/m/server/slave/internal/db"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_slave.RegisterSlaveServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	db.Init()
	return grpc_slave.RegisterSlaveHandler(context.Background(), gateway, cc)
}
