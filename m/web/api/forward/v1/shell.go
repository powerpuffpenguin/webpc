package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Shell struct {
	web.Helper
}

func (h Shell) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`shell`)
	r.GET(`:id/:shellid/:cols/:rows`, h.connect)
}
func (h Shell) connect(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()
	var obj struct {
		ID      string `uri:"id" binding:"required"`
		Shellid int64  `uri:"shellid" `
		Cols    uint32 `uri:"cols" binding:"required"`
		Rows    uint32 `uri:"rows" binding:"required"`
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
	client := grpc_shell.NewShellClient(cc)
	stream, e := client.Connect(ctx)
	if e != nil {
		ws.Error(e)
		return
	}
	e = stream.Send(&grpc_shell.ConnectRequest{
		Event: grpc_shell.Event_Connect,
		Id:    obj.Shellid,
		Cols:  obj.Cols,
		Rows:  obj.Rows,
	})
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		if messageType == websocket.BinaryMessage {
			return stream.Send(&grpc_shell.ConnectRequest{
				Event:  grpc_shell.Event_Binary,
				Binary: p,
			})
		}
		var req grpc_shell.ConnectRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return stream.Send(&req)
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		} else if resp.Event == grpc_shell.Event_Binary {
			return ws.SendBinary(resp.Binary)
		}
		return ws.SendMessage(resp)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
