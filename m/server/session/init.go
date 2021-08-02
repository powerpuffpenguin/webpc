package session

import (
	"context"
	grpc_session "github.com/powerpuffpenguin/webpc/protocol/session"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_session.RegisterSessionServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_session.RegisterSessionHandler(context.Background(), gateway, cc)
}
