package forward

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	grpc_forward "github.com/powerpuffpenguin/webpc/protocol/forward/forward"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_forward.RegisterForwardServer(srv, server{})

	if ce := logger.Logger.Check(zap.InfoLevel, `port forward`); ce != nil {
		ce.Write(
			zap.Bool(`enable`, configure.DefaultSystem().PortForward),
		)
	}
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	return nil
}
