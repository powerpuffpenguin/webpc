package db

import (
	"context"
	"strconv"
	"time"

	"github.com/powerpuffpenguin/webpc/db/manipulator"
	grpc_user "github.com/powerpuffpenguin/webpc/protocol/user"
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Find(ctx context.Context, request *grpc_user.FindRequest) (*grpc_user.FindResponse, error) {
	if request.Result < grpc_user.FindRequest_DATA ||
		request.Result > grpc_user.FindRequest_DATA_COUNT {
		return nil, status.Error(codes.InvalidArgument, `not support result enum : `+strconv.FormatInt(int64(request.Result), 10))
	}

	session := manipulator.Session(ctx)
	defer session.Close()
	if request.Parent > 1 {
		resp, e := signal_group.IDS(ctx, request.Parent, true)
		if e != nil {
			return nil, e
		} else if len(resp.Args) == 0 {
			return &grpc_user.FindResponse{
				Result: request.Result,
			}, nil
		}
		session.In(colParent, resp.Args...)
	}

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
func Add(ctx context.Context, parent int64, name, nickname, password string, authorization []int32) (id int64, e error) {
	bean := &DataOfUser{
		Parent:        parent,
		Name:          name,
		Nickname:      nickname,
		Password:      password,
		Authorization: authorization,
	}
	resp, e := signal_group.Exists(ctx, parent)
	if e != nil {
		return 0, e
	} else if !resp.Exists {
		return 0, status.Error(codes.NotFound, `group not exists: `+strconv.FormatInt(parent, 10))
	}

	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return
	}
	defer session.Close()

	// exec
	_, e = session.InsertOne(bean)
	if e != nil {
		return
	}
	_, e = modtimeHelper.Modified(session, time.Now())
	if e != nil {
		return
	}

	// commit
	e = session.Commit()
	if e != nil {
		return
	}
	id = bean.ID
	return
}

// Password change password
func Password(ctx context.Context, id int64, value string) (bool, error) {
	rowsAffected, e := manipulator.Engine().
		ID(id).
		Context(ctx).
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
func Change(ctx context.Context, id int64, nickname string, authorization []int32) (bool, error) {
	// begin
	session, e := manipulator.Begin(ctx)
	if e != nil {
		return false, e
	}
	defer session.Close()

	// exec
	rowsAffected, e := session.
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
		Update(&DataOfUser{
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
		In(colID, args...).
		Delete(&DataOfUser{})
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
