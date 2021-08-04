package db

import (
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	signal_slave "github.com/powerpuffpenguin/webpc/signal/slave"
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

	signal_slave.ConnectGroup(soltGroup)
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
func soltGroup(req *signal_slave.GroupRequest, resp *signal_slave.GroupResponse) (e error) {
	rowsAffected, e := req.Session.
		In(colParent, req.Args...).
		Cols(colParent).
		Update(&DataOfSlave{
			Parent: 0,
		})
	if e != nil {
		return
	}
	resp.RowsAffected = rowsAffected
	return
}
