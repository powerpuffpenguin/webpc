package db

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	grpc_slave "github.com/powerpuffpenguin/webpc/protocol/slave"

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

	session := manipulator.Session().Context(ctx)
	defer session.Close()
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
func Add(name, description string) (id int64, code string, e error) {
	code, e = newCode()
	if e != nil {
		return
	}

	bean := &DataOfSlave{
		Name:        name,
		Description: description,
		Code:        code,
	}
	_, e = manipulator.Engine().InsertOne(bean)
	if e != nil {
		return
	}
	id = bean.ID
	modtimeHelper.Modified(time.Now())
	return
}

func Code(id int64) (changed bool, code string, e error) {
	code, e = newCode()
	if e != nil {
		return
	}

	rowsAffected, e := manipulator.Engine().
		ID(id).
		Cols(colCode).
		Update(&DataOfSlave{
			Code: code,
		})
	if e != nil {
		return
	}
	changed = rowsAffected != 0
	return
}

// Change properties
func Change(id int64, name, description string) (bool, error) {
	rowsAffected, e := manipulator.Engine().
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
		modtimeHelper.Modified(time.Now())
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
	rowsAffected, e := manipulator.Engine().
		Context(ctx).
		In(colID, args...).
		Delete(&DataOfSlave{})
	if e != nil {
		return 0, e
	}
	if rowsAffected != 0 {
		modtimeHelper.Modified(time.Now())
	}
	return rowsAffected, nil
}
