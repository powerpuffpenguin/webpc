package shell

import (
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
)

func Init() {
	cnf := configure.DefaultSystem()
	command = cnf.Shell
	if ce := logger.Logger.Check(zap.InfoLevel, `shell`); ce != nil {
		ce.Write(
			zap.String("shell", command),
		)
	}

	DefaultManager().restore()
}
