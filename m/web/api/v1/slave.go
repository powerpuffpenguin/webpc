package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
	"google.golang.org/grpc"
)

type Slave struct {
	web.Helper
	cc *grpc.ClientConn
}

func (h *Slave) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	h.cc = cc
	r := router.Group(`slaves`)
	r.GET(`subscribe`, h.subscribe)
}

func (h *Slave) subscribe(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()

	ctx := h.NewContext(c)
	// fmt.Println(`--------------- new subscribe`)
	// defer fmt.Println(`--------------- exit subscribe`)
	client := grpc_slave.NewSlaveClient(h.cc)
	stream, e := client.Subscribe(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		var req grpc_slave.SubscribeRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
