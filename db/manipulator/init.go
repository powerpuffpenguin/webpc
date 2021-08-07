package manipulator

import (
	"context"
	"fmt"

	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/db/data"
	"github.com/powerpuffpenguin/webpc/logger"

	"go.uber.org/zap"
	"xorm.io/xorm"
	"xorm.io/xorm/caches"
	"xorm.io/xorm/schemas"
)

var _Engine *xorm.EngineGroup

// Init Initialize the database
func Init(cnf *configure.DB) {
	engine, e := xorm.NewEngineGroup(cnf.Driver, cnf.Source)
	if e != nil {
		logger.Logger.Panic(`new db engine group error`,
			zap.Error(e),
			zap.String(`driver`, cnf.Driver),
			zap.Strings(`source`, cnf.Source),
		)
	}
	_Engine = engine
	if ce := logger.Logger.Check(zap.InfoLevel, `db group`); ce != nil {
		ce.Write(
			zap.String(`driver`, cnf.Driver),
			zap.Strings(`source`, cnf.Source),
		)
	}
	engine.ShowSQL(cnf.ShowSQL)
	// cache
	if cnf.Cache.Record > 0 {
		engine.SetDefaultCacher(
			caches.NewLRUCacher(
				caches.NewMemoryStore(),
				cnf.Cache.Record,
			),
		)
		if ce := logger.Logger.Check(zap.InfoLevel, `db default cacher`); ce != nil {
			ce.Write(
				zap.Int(`record`, cnf.Cache.Record),
			)
		}
		for _, name := range cnf.Cache.Direct {
			engine.SetCacher(name, nil)
			if ce := logger.Logger.Check(zap.InfoLevel, `db disable cacher`); ce != nil {
				ce.Write(
					zap.String(`table name`, name),
				)
			}
		}
	}
	for _, item := range cnf.Cache.Special {
		if item.Name != `` && item.Record > 0 {
			engine.SetCacher(item.Name, caches.NewLRUCacher(
				caches.NewMemoryStore(),
				item.Record,
			))
			if ce := logger.Logger.Check(zap.InfoLevel, `db set cacher`); ce != nil {
				ce.Write(
					zap.String(`name`, item.Name),
					zap.Int(`record`, item.Record),
				)
			}
		}
	}
	// connect pool
	if cnf.MaxOpen > 1 {
		engine.SetMaxOpenConns(cnf.MaxOpen)
	}
	if cnf.MaxIdle > 1 {
		engine.SetMaxIdleConns(cnf.MaxIdle)
	}
	// table
	initTable(engine)
}

func initTable(engine *xorm.EngineGroup) {
	session, e := Begin()
	if e != nil {
		logger.Logger.Panic(`db begin error`,
			zap.Error(e),
		)
	}
	defer session.Close()
	e = SyncTable(session,
		&data.Modtime{},
	)
	if e != nil {
		logger.Logger.Panic(`sync table error`,
			zap.Error(e),
		)
	}
	e = session.Commit()
	if e != nil {
		logger.Logger.Panic(`sync table error`,
			zap.Error(e),
		)
	}
}

// SyncTable sync table to consistent
func SyncTable(session *xorm.Session, beans ...interface{}) (e error) {
	for i := 0; i < len(beans); i++ {
		e = syncTable(session, beans[i])
		if e != nil {
			return
		}
	}
	return
}
func syncTable(session *xorm.Session, bean interface{}) (e error) {
	has, e := session.IsTableExist(bean)
	if e != nil {
		if ce := logger.Logger.Check(zap.PanicLevel, `IsTableExist`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	} else if has {
		e = session.Sync2(bean)
		if e != nil {
			if ce := logger.Logger.Check(zap.PanicLevel, `Sync2`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			return
		}
	} else {
		e = session.CreateTable(bean)
		if e != nil {
			if ce := logger.Logger.Check(zap.PanicLevel, `CreateTable`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			return
		}
		e = session.CreateIndexes(bean)
		if e != nil {
			if ce := logger.Logger.Check(zap.PanicLevel, `CreateIndexes`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			return
		}
		e = session.CreateUniques(bean)
		if e != nil {
			if ce := logger.Logger.Check(zap.PanicLevel, `CreateUniques`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			return
		}
	}
	return
}

func Engine() *xorm.EngineGroup {
	return _Engine
}
func Master() *xorm.Engine {
	return _Engine.Master()
}
func Slave() *xorm.Engine {
	return _Engine.Slave()
}
func Session(ctx ...context.Context) *xorm.Session {
	s := _Engine.NewSession()
	for _, c := range ctx {
		s = s.Context(c)
	}
	return s
}
func Begin(ctx ...context.Context) (*xorm.Session, error) {
	s := Session(ctx...)
	e := s.Begin()
	if e != nil {
		s.Close()
		return nil, e
	}
	return s, nil
}
func DeleteCache(tableName string, pk ...interface{}) error {
	if len(pk) == 0 {
		return nil
	}
	cacher := _Engine.GetCacher(tableName)
	if cacher == nil {
		return nil
	}
	id, e := schemas.NewPK(pk...).ToString()
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `DeleteCache`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`table name`, tableName),
				zap.String(`pk`, fmt.Sprint(pk...)),
			)
		}
		return e
	}
	cacher.ClearIds(tableName)
	cacher.DelBean(tableName, id)
	return nil
}
