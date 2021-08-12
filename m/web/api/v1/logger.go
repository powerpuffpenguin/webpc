package v1

import (
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_logger "github.com/powerpuffpenguin/webpc/protocol/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type Logger struct {
	web.Helper
	cc *grpc.ClientConn
}

func (h *Logger) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	h.cc = cc
	r := router.Group(`logger`)
	r.GET(`attach`, h.attach)
}

func (h *Logger) attach(c *gin.Context) {
	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	defer ws.Close()
	ctx := h.NewContext(c)
	// fmt.Println(`--------------- new attach`)
	// defer fmt.Println(`--------------- exit attach`)
	client := grpc_logger.NewLoggerClient(h.cc)
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
