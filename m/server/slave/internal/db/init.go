package db

import (
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
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

	signal_group.ConnectDelete(soltGroupDelete)
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
func soltGroupDelete(req *signal_group.DeleteRequest, resp *signal_group.DeleteResponse) (e error) {
	rowsAffected, e := req.Session.
		In(colParent, req.Args...).
		Cols(colParent).
		Update(&DataOfSlave{
			Parent: 0,
		})
	if e != nil {
		return
	}
	resp.RowsAffected += rowsAffected
	if rowsAffected != 0 {
		modtimeHelper.Modified(req.Session, time.Now())
	}
	return
}
