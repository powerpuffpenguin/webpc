package web

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type IHelper interface {
	Register(*grpc.ClientConn, *gin.RouterGroup)
}
