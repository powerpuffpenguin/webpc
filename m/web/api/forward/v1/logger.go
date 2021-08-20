package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_logger "github.com/powerpuffpenguin/webpc/protocol/forward/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Logger struct {
	web.Helper
}

func (h Logger) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`logger`)
	r.GET(`attach/:id`, h.attach)
}
func (h Logger) attach(c *gin.Context) {
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

	cc, e := forward.Default().Get(obj.ID)
	if e != nil {
		ws.Error(e)
		return
	}

	ctx := h.NewForwardContext(c)
	client := grpc_logger.NewLoggerClient(cc)
	stream, e := client.Attach(ctx, &grpc_logger.AttachRequest{})
	if e != nil {
		ws.Error(e)
		return
	}
	f := web.NewForward(func(counted uint64, messageType int, p []byte) error {
		return nil
	}, func(counted uint64) (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		if counted == 0 && len(resp.Data) == 0 {
			return ws.Success()
		}
		return ws.SendBinary(resp.Data)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
