package slave

import (
	"context"
	"errors"
)

type CodeRequest struct {
	Context context.Context
	Code    string
}

type CodeResponse struct {
	ID     int64
	Parent int64
}
type CodeHandle func(req *CodeRequest, resp *CodeResponse) (e error)

var defaultCode []CodeHandle

func ConnectCode(f CodeHandle) {
	defaultCode = append(defaultCode, f)
}
func Code(ctx context.Context,
	code string,
) (result *CodeResponse, e error) {
	if len(defaultCode) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &CodeRequest{
		Context: ctx,
		Code:    code,
	}
	var resp CodeResponse
	for _, f := range defaultCode {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
