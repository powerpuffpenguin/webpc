package db

import (
	"context"
	"strconv"
	"github.com/powerpuffpenguin/webpc/db/manipulator"
	grpc_user "github.com/powerpuffpenguin/webpc/protocol/user"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Find(ctx context.Context, request *grpc_user.FindRequest) (*grpc_user.FindResponse, error) {
	if request.Result < grpc_user.FindRequest_DATA ||
		request.Result > grpc_user.FindRequest_DATA_COUNT {
		return nil, status.Error(codes.InvalidArgument, `not support result enum : `+strconv.FormatInt(int64(request.Result), 10))
	}

	session := manipulator.Session().Context(ctx)
	defer session.Close()
	var beans []DataOfUser
	var response grpc_user.FindResponse
	response.Result = request.Result
	if request.Name != `` {
		if request.NameFuzzy {
			session.Where(colName+` like ?`, `%`+request.Name+`%`)
		} else {
			session.Where(colName+` = ?`, request.Name)
		}
	}
	switch request.Result {
	case grpc_user.FindRequest_COUNT:
		count, e := session.Count(&DataOfUser{})
		if e != nil {
			return nil, e
		}
		response.Count = count
	case grpc_user.FindRequest_DATA:
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
		response.Data = make([]*grpc_user.Data, len(beans))
		for i := 0; i < len(beans); i++ {
			response.Data[i] = beans[i].ToPB()
		}
	}
	return &response, nil
}
func Add(name, nickname, password string, authorization []int32) (int64, error) {
	bean := &DataOfUser{
		Name:          name,
		Nickname:      nickname,
		Password:      password,
		Authorization: authorization,
	}
	_, e := manipulator.Engine().InsertOne(bean)
	if e != nil {
		return 0, e
	}
	modtimeHelper.Modified(time.Now())
	return bean.ID, nil
}

// Password change password
func Password(id int64, value string) (bool, error) {
	rowsAffected, e := manipulator.Engine().
		ID(id).
		Cols(colPassword).
		Update(&DataOfUser{
			Password: value,
		})
	if e != nil {
		return false, e
	}
	return rowsAffected != 0, nil
}

// Change properties
func Change(id int64, nickname string, authorization []int32) (bool, error) {
	rowsAffected, e := manipulator.Engine().
		ID(id).
		Cols(colNickname, colAuthorization).
		Update(&DataOfUser{
			Nickname:      nickname,
			Authorization: authorization,
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
		Delete(&DataOfUser{})
	if e != nil {
		return 0, e
	}
	if rowsAffected != 0 {
		modtimeHelper.Modified(time.Now())
	}
	return rowsAffected, nil
}
