package slave

import (
	"context"
	"errors"

	"xorm.io/xorm"
)

type DeleteRequest struct {
	Context context.Context
	Session *xorm.Session
	Args    []interface{}
}

type DeleteResponse struct {
	RowsAffected int64
}
type DeleteHandle func(req *DeleteRequest, resp *DeleteResponse) (e error)

var defaultDelete []DeleteHandle

func ConnectDelete(f DeleteHandle) {
	defaultDelete = append(defaultDelete, f)
}
func Delete(ctx context.Context,
	session *xorm.Session,
	args []interface{},
) (result *DeleteResponse, e error) {
	if len(defaultDelete) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &DeleteRequest{
		Context: ctx,
		Session: session,
		Args:    args,
	}
	var resp DeleteResponse
	for _, f := range defaultDelete {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
