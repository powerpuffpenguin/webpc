package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/web"
	"github.com/powerpuffpenguin/webpc/m/web/api/v1/internal/dialer"
	signal_slave "github.com/powerpuffpenguin/webpc/signal/slave"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Dialer struct {
	web.Helper
}

func (h Dialer) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`dialer`)
	r.GET(`:code`, h.CheckWebsocket, h.dialer)
}
func (h Dialer) dialer(c *gin.Context) {
	var obj struct {
		Code string `uri:"code" binding:"required"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	ctx := c.Request.Context()
	resp, e := signal_slave.Code(ctx, obj.Code)
	if e != nil {
		h.Error(c, e)
		return
	} else if resp.ID == 0 {
		h.Error(c, status.Error(codes.NotFound, `code not exists: `+obj.Code))
		return
	}

	ws, e := h.Websocket(c, nil)
	if e != nil {
		return
	}
	conn := dialer.NewConn(resp.ID, ws.UnderlyingConn())
	dialer.Put(conn)
	<-conn.Done()
}
