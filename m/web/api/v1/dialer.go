package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/web"
	"github.com/powerpuffpenguin/webpc/m/web/api/v1/internal/dialer"
	"google.golang.org/grpc"
)

type Dialer struct {
	web.Helper
}

func (h Dialer) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`dialer`)
	r.GET(``, h.CheckWebsocket, h.dialer)
}
func (h Dialer) dialer(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		return
	}
	conn := dialer.NewConn(1, ws.UnderlyingConn())
	fmt.Println(conn.RemoteAddr())
	dialer.Put(conn)
	<-conn.Done()
}
