package group

import (
	"context"
	"errors"
)

type IDSRequest struct {
	Context context.Context
	ID      int64
	Args    bool
}

type IDSResponse struct {
	ID   []int64
	Args []interface{}
}
type IDSHandle func(req *IDSRequest, resp *IDSResponse) (e error)

var defaultIDS []IDSHandle

func ConnectIDS(f IDSHandle) {
	defaultIDS = append(defaultIDS, f)
}
func IDS(ctx context.Context,
	id int64,
	args bool,
) (result *IDSResponse, e error) {
	if len(defaultIDS) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &IDSRequest{
		Context: ctx,
		ID:      id,
		Args:    args,
	}
	var resp IDSResponse
	for _, f := range defaultIDS {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
