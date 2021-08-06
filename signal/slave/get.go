package slave

import (
	"context"
	"errors"
)

type GetRequest struct {
	Context context.Context
	ID      int64
}

type GetResponse struct {
	ID          int64
	Name        string
	Description string
	Code        string
	Parent      int64
}
type GetHandle func(req *GetRequest, resp *GetResponse) (e error)

var defaultGet []GetHandle

func ConnectGet(f GetHandle) {
	defaultGet = append(defaultGet, f)
}
func Get(ctx context.Context,
	id int64,
) (result *GetResponse, e error) {
	if len(defaultGet) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &GetRequest{
		Context: ctx,
		ID:      id,
	}
	var resp GetResponse
	for _, f := range defaultGet {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
