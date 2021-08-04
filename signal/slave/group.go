package slave

import (
	"context"
	"errors"

	"xorm.io/xorm"
)

type GroupRequest struct {
	Context context.Context
	Session *xorm.Session
	Args    []interface{}
}

type GroupResponse struct {
	RowsAffected int64
}
type GroupHandle func(req *GroupRequest, resp *GroupResponse) (e error)

var defaultGroup []GroupHandle

func ConnectGroup(f GroupHandle) {
	defaultGroup = append(defaultGroup, f)
}
func Group(ctx context.Context,
	session *xorm.Session,
	args []interface{},
) (result *GroupResponse, e error) {
	if len(defaultGroup) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &GroupRequest{
		Context: ctx,
		Session: session,
		Args:    args,
	}
	var resp GroupResponse
	for _, f := range defaultGroup {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
