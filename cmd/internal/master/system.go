package master

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/m/register"
	"github.com/powerpuffpenguin/webpc/slave"
	"google.golang.org/grpc"
)

func registerSystem(srv *grpc.Server, system *configure.System, cc *grpc.ClientConn, gateway *runtime.ServeMux) (e error) {
	if !system.Enable {
		return
	}
	e = slave.HTTP(gateway, cc)
	if e != nil {
		return
	}
	slave.GRPC(srv)

	register.DefaultForward().Put(0, cc, gateway)
	return
}
