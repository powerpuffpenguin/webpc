package group

import (
	"context"
	"errors"
)

type ExistsRequest struct {
	Context context.Context
	ID      int64
}

type ExistsResponse struct {
	Exists bool
}
type ExistsHandle func(req *ExistsRequest, resp *ExistsResponse) (e error)

var defaultExists []ExistsHandle

func ConnectExists(f ExistsHandle) {
	defaultExists = append(defaultExists, f)
}
func Exists(ctx context.Context,
	id int64,
) (result *ExistsResponse, e error) {
	if len(defaultExists) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &ExistsRequest{
		Context: ctx,
		ID:      id,
	}
	var resp ExistsResponse
	for _, f := range defaultExists {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
