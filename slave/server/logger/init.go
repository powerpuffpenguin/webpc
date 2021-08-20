package logger

import (
	"context"

	grpc_logger "github.com/powerpuffpenguin/webpc/protocol/forward/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_logger.RegisterLoggerServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return grpc_logger.RegisterLoggerHandler(context.Background(), gateway, cc)
}
