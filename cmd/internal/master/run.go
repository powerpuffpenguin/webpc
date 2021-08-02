package master

import (
	"net"

	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run(cnf *configure.HTTP, debug bool) {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	// listen
	l, e := net.Listen(`tcp`, cnf.Addr)
	if e != nil {
		logger.Logger.Panic(`listen error`,
			zap.Error(e),
		)
		return
	}
	h2 := cnf.H2()
	logger.Logger.Info(`listen success`,
		zap.String(`addr`, cnf.Addr),
		zap.Bool(`h2`, h2),
	)
	// serve
	s := newServer(l, cnf.Swagger, debug, &cnf.Option)
	if h2 {
		e = s.ServeTLS(cnf.CertFile, cnf.KeyFile)
	} else {
		e = s.Serve()
	}
	if e != nil {
		logger.Logger.Panic(`serve error`,
			zap.Error(e),
		)
		return
	}
}
