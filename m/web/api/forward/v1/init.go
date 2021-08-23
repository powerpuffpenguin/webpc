package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/powerpuffpenguin/webpc/m/web"
	"google.golang.org/grpc"
)

const BaseURL = `v1`

type Helper struct {
	web.Helper
}

func (h Helper) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		Filesystem{},
		Static{},
		Logger{},
		Shell{},
		VNC{},
	}
	for _, m := range ms {
		m.Register(cc, r)
	}
}
