package mount

import (
	"errors"

	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
)

func Init(ms []configure.Mount) {
	fs := Default()
	keys := make(map[string]bool)
	for _, item := range ms {
		if keys[item.Name] {
			e := errors.New(`name already exists`)
			if ce := logger.Logger.Check(zap.FatalLevel, `mount error`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.String(`name`, item.Name),
					zap.String(`root`, item.Root),
					zap.Bool(`read`, item.Read),
					zap.Bool(`write`, item.Write),
					zap.Bool(`shared`, item.Shared),
				)
			}
			return
		}
		fs.Push(item.Name, item.Root, item.Read, item.Write, item.Shared)
		if ce := logger.Logger.Check(zap.InfoLevel, `mount`); ce != nil {
			ce.Write(
				zap.String(`name`, item.Name),
				zap.String(`root`, item.Root),
				zap.Bool(`read`, item.Read),
				zap.Bool(`write`, item.Write),
				zap.Bool(`shared`, item.Shared),
			)
		}
	}
}
