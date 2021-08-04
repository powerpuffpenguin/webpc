package db

import (
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	signal_slave "github.com/powerpuffpenguin/webpc/signal/slave"
	"golang.org/x/net/context"
)

func Add(ctx context.Context, pid int64, name, description string) (int64, time.Time, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return 0, time.Time{}, e
	}
	defer session.Close()

	// exec
	bean := &DataOfGroup{
		Parent:      pid,
		Name:        name,
		Description: description,
	}
	_, e = session.InsertOne(bean)
	if e != nil {
		return 0, time.Time{}, e
	}
	at := time.Now()
	_, e = modtimeHelper.Modified(session, at)
	if e != nil {
		return 0, time.Time{}, e
	}

	// commit
	e = session.Commit()
	if e != nil {
		return 0, time.Time{}, e
	}
	return bean.ID, at, nil
}
func Move(ctx context.Context, id, pid int64) (time.Time, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return time.Time{}, e
	}
	defer session.Close()

	// exec
	_, e = session.
		ID(id).
		Cols(colParent).
		Update(&DataOfGroup{
			Parent: pid,
		})
	if e != nil {
		return time.Time{}, e
	}
	at := time.Now()
	_, e = modtimeHelper.Modified(session, at)
	if e != nil {
		return time.Time{}, e
	}

	// commit
	e = session.Commit()
	if e != nil {
		return time.Time{}, e
	}
	return at, nil
}
func Change(ctx context.Context, id int64, name, description string) (time.Time, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return time.Time{}, e
	}
	defer session.Close()

	// exec
	_, e = session.
		ID(id).
		Cols(colName).
		Cols(colDescription).
		Update(&DataOfGroup{
			Name:        name,
			Description: description,
		})
	if e != nil {
		return time.Time{}, e
	}

	at := time.Now()
	_, e = modtimeHelper.Modified(session, at)
	if e != nil {
		return time.Time{}, e
	}

	// commit
	e = session.Commit()
	if e != nil {
		return time.Time{}, e
	}
	return at, nil
}
func Remove(ctx context.Context, args []interface{}) (time.Time, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return time.Time{}, e
	}
	defer session.Close()

	// exec
	_, e = session.
		In(colID, args...).
		Delete(&DataOfGroup{})
	if e != nil {
		return time.Time{}, e
	}
	_, e = signal_slave.Group(ctx, session, args)
	if e != nil {
		return time.Time{}, e
	}

	at := time.Now()
	_, e = modtimeHelper.Modified(session, at)
	if e != nil {
		return time.Time{}, e
	}

	// commit
	e = session.Commit()
	if e != nil {
		return time.Time{}, e
	}
	return at, nil
}
