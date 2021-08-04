package db

import (
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
)

var modtimeHelper = manipulator.ModtimeHelper(manipulator.ModtimeSlave)

func LastModified() (modtime time.Time) {
	modtime, _ = modtimeHelper.LastModified()
	return
}
func Init() {
	e := doInit()
	if e != nil {
		logger.Logger.Panic(`init db error`,
			zap.Error(e),
			zap.String(`table`, tableName),
		)
		return
	}
}
func doInit() (e error) {
	session, e := manipulator.Begin()
	if e != nil {
		return
	}
	defer session.Close()
	// modtime
	e = modtimeHelper.Init(session, time.Now().Unix(), ``, `slave`)
	if e != nil {
		return
	}
	// sync
	bean := &DataOfSlave{}
	e = manipulator.SyncTable(session, bean)
	if e != nil {
		return
	}
	e = session.Commit()
	return
}
