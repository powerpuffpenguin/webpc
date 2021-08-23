package vnc

import (
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_vnc "github.com/powerpuffpenguin/webpc/protocol/forward/vnc"
	"github.com/powerpuffpenguin/webpc/slave/server/vnc/internal/connect"
)

type server struct {
	grpc_vnc.UnimplementedVncServer
	helper.Helper
}

func (s server) Connect(server grpc_vnc.Vnc_ConnectServer) (e error) {
	w := connect.New(server)
	e = w.Serve()
	return
}
