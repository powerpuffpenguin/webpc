package v1

import (
	"errors"
	"net/http"

	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/m/web"
	grpc_logger "github.com/powerpuffpenguin/webpc/protocol/logger"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/sessionid"
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
func (h *Logger) checkRoot(c *gin.Context) (code int, msg string) {
	userdata, e := h.ShouldBindUserdata(c)
	if e != nil {
		if sessionid.IsTokenExpired(e) {
			code = http.StatusUnauthorized
		} else if errors.Is(e, sessionid.ErrTokenNotExists) {
			code = http.StatusForbidden
		} else {
			code = http.StatusInternalServerError
		}
		msg = e.Error()
		return
	}
	if !userdata.AuthAny(db.Root) {
		code = http.StatusForbidden
		msg = `permission denied`
		return
	}
	return
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
	e = ws.Success()
	if e != nil {
		return
	}
	f := web.NewForward(func(messageType int, p []byte) error {
		var req grpc_slave.SubscribeRequest
		e = web.Unmarshal(p, &req)
		if e != nil {
			return e
		}
		return nil
	}, func() (e error) {
		resp, e := stream.Recv()
		if e != nil {
			return
		}
		return ws.SendBinary(resp.Data)
	}, func() error {
		return stream.CloseSend()
	})
	ws.Forward(f)
}
