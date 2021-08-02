package v1

import (
	"github.com/powerpuffpenguin/webpc/m/web"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

const BaseURL = `v1`

type Helper struct {
	web.Helper
}

func (h Helper) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		Logger{},
	}
	for _, m := range ms {
		m.Register(cc, r)
	}
}
