package forward

import (
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_forward "github.com/powerpuffpenguin/webpc/protocol/forward/forward"
	"github.com/powerpuffpenguin/webpc/slave/server/forward/internal/connect"
)

type server struct {
	grpc_forward.UnimplementedForwardServer
	helper.Helper
}

func (s server) Connect(server grpc_forward.Forward_ConnectServer) (e error) {
	w := connect.New(server)
	e = w.Serve()
	return
}
