package forward

import (
	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/web"
	v1 "github.com/powerpuffpenguin/webpc/m/web/api/forward/v1"
	"google.golang.org/grpc"
)

const BaseURL = `forward`

type Helper struct {
	web.Helper
}

func (h Helper) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		v1.Helper{},
	}
	for _, m := range ms {
		m.Register(cc, r)
	}
}
