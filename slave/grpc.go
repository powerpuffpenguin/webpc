package slave

import (
	m_fs "github.com/powerpuffpenguin/webpc/slave/server/fs"
	m_system "github.com/powerpuffpenguin/webpc/slave/server/system"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var ms = []Module{
	m_fs.Module(0),
	m_system.Module(0),
}

func GRPC(srv *grpc.Server) {
	for _, m := range ms {
		m.RegisterGRPC(srv)
	}
}
func HTTP(gateway *runtime.ServeMux, cc *grpc.ClientConn) (e error) {
	for _, m := range ms {
		e = m.RegisterGateway(gateway, cc)
		if e != nil {
			return
		}
	}
	return
}

type Module interface {
	RegisterGRPC(srv *grpc.Server)
	RegisterGateway(gateway *runtime.ServeMux, cc *grpc.ClientConn) error
}
