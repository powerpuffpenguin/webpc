package vnc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	grpc_vnc "github.com/powerpuffpenguin/webpc/protocol/forward/vnc"
	"github.com/powerpuffpenguin/webpc/slave/server/vnc/internal/connect"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Module int

func (Module) RegisterGRPC(srv *grpc.Server) {
	grpc_vnc.RegisterVncServer(srv, server{})
}
func (Module) RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error {
	vnc := configure.DefaultSystem().VNC
	connect.VNC = vnc
	if ce := logger.Logger.Check(zap.InfoLevel, `VNC`); ce != nil {
		ce.Write(zap.String(`connect`, vnc))
	}
	return nil
}
