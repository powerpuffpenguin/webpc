package db

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newCode() (val string, e error) {
	u, e := uuid.NewUUID()
	if e != nil {
		return
	}
	val = hex.EncodeToString(u[:])
	return
}
func Find(ctx context.Context, request *grpc_slave.FindRequest) (*grpc_slave.FindResponse, error) {
	if request.Result < grpc_slave.FindRequest_DATA ||
		request.Result > grpc_slave.FindRequest_DATA_COUNT {
		return nil, status.Error(codes.InvalidArgument, `not support result enum : `+strconv.FormatInt(int64(request.Result), 10))
	}

	session := manipulator.Session(ctx)
	defer session.Close()
	if request.Parent > 1 {
		resp, e := signal_group.IDS(ctx, request.Parent, true)
		if e != nil {
			return nil, e
		} else if len(resp.Args) == 0 {
			return &grpc_slave.FindResponse{
				Result: request.Result,
			}, nil
		}
		session.In(colParent, resp.Args...)
	}
	var beans []DataOfSlave
	var response grpc_slave.FindResponse
	response.Result = request.Result
	if request.Name != `` {
		if request.NameFuzzy {
			session.Where(colName+` like ?`, `%`+request.Name+`%`)
		} else {
			session.Where(colName+` = ?`, request.Name)
		}
	}
	switch request.Result {
	case grpc_slave.FindRequest_COUNT:
		count, e := session.Count(&DataOfSlave{})
		if e != nil {
			return nil, e
		}
		response.Count = count
	case grpc_slave.FindRequest_DATA:
		e := session.OrderBy(colID).Limit(int(request.Limit), int(request.Offset)).Find(&beans)
		if e != nil {
			return nil, e
		}
	default: // FindRequest_DATA_COUNT
		count, e := session.OrderBy(colID).Limit(int(request.Limit), int(request.Offset)).FindAndCount(&beans)
		if e != nil {
			return nil, e
		}
		response.Count = count
	}
	if len(beans) != 0 {
		response.Data = make([]*grpc_slave.Data, len(beans))
		for i := 0; i < len(beans); i++ {
			response.Data[i] = beans[i].ToPB()
		}
	}
	return &response, nil
}
func Get(ctx context.Context, id int64) (*grpc_slave.Data, error) {
	var bean DataOfSlave
	exists, e := manipulator.Engine().ID(id).Context(ctx).Get(&bean)
	if e != nil {
		return nil, e
	} else if !exists {
		return nil, status.Error(codes.NotFound, `id not exists: `+strconv.FormatInt(id, 10))
	}
	return bean.ToPB(), nil
}
func Add(ctx context.Context, parent int64, name, description string) (int64, string, error) {
	code, e := newCode()
	if e != nil {
		return 0, ``, e
	}
	bean := &DataOfSlave{
		Parent:      parent,
		Name:        name,
		Description: description,
		Code:        code,
	}
	resp, e := signal_group.Exists(ctx, parent)
	if e != nil {
		return 0, ``, e
	} else if !resp.Exists {
		return 0, ``, status.Error(codes.NotFound, `group not exists: `+strconv.FormatInt(parent, 10))
	}

	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return 0, ``, e
	}
	defer session.Close()

	// exec
	_, e = session.InsertOne(bean)
	if e != nil {
		return 0, ``, e
	}
	_, e = modtimeHelper.Modified(session, time.Now())
	if e != nil {
		return 0, ``, e
	}

	// commit
	e = session.Commit()
	if e != nil {
		return 0, ``, e
	}
	return bean.ID, code, nil
}

func Code(ctx context.Context, id int64) (bool, string, error) {
	code, e := newCode()
	if e != nil {
		return false, ``, e
	}

	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return false, ``, e
	}
	defer session.Close()

	// exec
	rowsAffected, e := session.
		ID(id).
		Cols(colCode).
		Update(&DataOfSlave{
			Code: code,
		})
	if e != nil {
		return false, ``, e
	}
	changed := rowsAffected != 0
	if changed {
		_, e = modtimeHelper.Modified(session, time.Now())
		if e != nil {
			return false, ``, e
		}
	}

	// commit
	e = session.Commit()
	if e != nil {
		return false, ``, e
	}
	return changed, code, nil
}

// Change properties
func Change(ctx context.Context, id int64, name, description string) (bool, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return false, e
	}
	defer session.Close()

	// exec
	rowsAffected, e := session.
		ID(id).
		Cols(colName, colDescription).
		Update(&DataOfSlave{
			Name:        name,
			Description: description,
		})
	if e != nil {
		return false, e
	}
	changed := rowsAffected != 0
	if changed {
		_, e = modtimeHelper.Modified(session, time.Now())
		if e != nil {
			return false, e
		}
	}

	// commit
	e = session.Commit()
	if e != nil {
		return false, e
	}
	return changed, nil
}
func Parent(ctx context.Context, id, parent int64) (bool, error) {
	resp, e := signal_group.Exists(ctx, parent)
	if e != nil {
		return false, e
	} else if !resp.Exists {
		return false, status.Error(codes.NotFound, `group not exists: `+strconv.FormatInt(parent, 10))
	}

	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return false, e
	}
	defer session.Close()

	// exec
	rowsAffected, e := session.
		ID(id).
		Cols(colParent).
		Update(&DataOfSlave{
			Parent: parent,
		})
	if e != nil {
		return false, e
	}
	changed := rowsAffected != 0
	if changed {
		_, e = modtimeHelper.Modified(session, time.Now())
		if e != nil {
			return false, e
		}
	}

	// commit
	e = session.Commit()
	if e != nil {
		return false, e
	}
	return changed, nil
}
func Remove(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return 0, e
	}
	defer session.Close()

	// exec
	rowsAffected, e := session.
		Context(ctx).
		In(colID, args...).
		Delete(&DataOfSlave{})
	if e != nil {
		return 0, e
	}
	changed := rowsAffected != 0
	if changed {
		_, e = modtimeHelper.Modified(session, time.Now())
		if e != nil {
			return 0, e
		}
	}

	// commit
	e = session.Commit()
	if e != nil {
		return 0, e
	}
	if changed {
		manipulator.DeleteCache(tableName, args...)
	}
	return rowsAffected, nil
}
