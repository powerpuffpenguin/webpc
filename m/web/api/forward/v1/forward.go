package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_forward "github.com/powerpuffpenguin/webpc/protocol/forward/forward"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Forward struct {
	web.Helper
}

func (h Forward) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`forward`)
	r.GET(`:id`, h.connect)
}
func (h Forward) connect(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()
	var obj struct {
		ID string `uri:"id" binding:"required"`
	}
	e = c.ShouldBindUri(&obj)
	if e != nil {
		ws.Error(status.Error(codes.InvalidArgument, e.Error()))
		return
	}
	ctx, cc, e := forward.Default().Get(c, obj.ID)
	if e != nil {
		ws.Error(e)
		return
	}
	client := grpc_forward.NewForwardClient(cc)
	stream, e := client.Connect(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	e = stream.Send(&grpc_forward.ConnectRequest{
		Event: grpc_forward.Event_Connect,
	})
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		if messageType == websocket.BinaryMessage {
			return stream.Send(&grpc_forward.ConnectRequest{
				Event:  grpc_forward.Event_Binary,
				Binary: p,
			})
		}
		var req grpc_forward.ConnectRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		} else if resp.Event == grpc_forward.Event_Binary {
			return ws.SendBinary(resp.Binary)
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
