package db

import (
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	"github.com/powerpuffpenguin/webpc/logger"

	"go.uber.org/zap"
)

const RootID = 1

var modtimeHelper = manipulator.ModtimeHelper(manipulator.ModtimeGroup)

func LastModified() (modtime time.Time) {
	modtime, _ = modtimeHelper.LastModified()
	return
}
func Init() map[int64]*DataOfGroup {
	// disabled cache
	manipulator.Engine().SetCacher(tableName, nil)

	// init db
	keys, e := doInit()
	if e != nil {
		logger.Logger.Panic(`init db error`,
			zap.Error(e),
			zap.String(`table`, tableName),
		)
		return nil
	}

	return keys
}
func doInit() (keys map[int64]*DataOfGroup, e error) {
	session, e := manipulator.Begin()
	if e != nil {
		return
	}
	defer session.Close()
	// modtime
	e = modtimeHelper.Init(session, time.Now().Unix(), ``, `group`)
	if e != nil {
		return
	}
	// sync
	bean := &DataOfGroup{}
	e = manipulator.SyncTable(session, bean)
	if e != nil {
		return
	}
	// find
	var items []DataOfGroup
	e = session.Find(&items)
	if e != nil {
		return
	}
	// keys
	keys = make(map[int64]*DataOfGroup)
	for i := range items {
		if items[i].ID == 0 {
			continue
		}
		keys[items[i].ID] = &items[i]
	}
	if ele, exists := keys[RootID]; exists {
		ele.Parent = 0
	} else {
		// insert root
		item := &DataOfGroup{
			ID:          RootID,
			Name:        `root`,
			Description: `include all devices`,
		}
		_, e = session.InsertOne(item)
		if e != nil {
			return
		}
		keys[RootID] = item
	}
	e = session.Commit()
	return
}
