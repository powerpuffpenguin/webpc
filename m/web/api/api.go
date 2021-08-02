package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/m/web"
	v1 "github.com/powerpuffpenguin/webpc/m/web/api/v1"
	"google.golang.org/grpc"
)

const BaseURL = `api`

type Helper struct {
	web.Helper
}

var maxBytesReader int64

func (h Helper) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	maxBytesReader = int64(configure.DefaultConfigure().HTTP.Option.MaxRecvMsgSize)
	r := router.Group(BaseURL)
	if maxBytesReader > 0 {
		r.Use(h.CheckBodySize)
	}

	ms := []web.IHelper{
		v1.Helper{},
	}
	for _, m := range ms {
		m.Register(cc, r)
	}
}

func (h Helper) CheckBodySize(c *gin.Context) {
	if maxBytesReader > 0 {
		if c.Request != nil && c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytesReader)
		}
	}
}
