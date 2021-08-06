package session

import (
	"context"
	"errors"
)

type SigninRequest struct {
	Context  context.Context
	Platform string
	Unix     int64
	Name     string
	Password string
}

type SigninResponse struct {
	ID int64
	// [a-zA-Z][a-zA-Z0-9]{4,}
	Name          string
	Nickname      string
	Authorization []int32
	Parent        int64
}
type SigninHandle func(req *SigninRequest, resp *SigninResponse) (e error)

var defaultSignin []SigninHandle

func ConnectSignin(f SigninHandle) {
	defaultSignin = append(defaultSignin, f)
}
func Signin(ctx context.Context,
	platform string, unix int64,
	name, password string,
) (result *SigninResponse, e error) {
	if len(defaultSignin) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &SigninRequest{
		Context:  ctx,
		Platform: platform,
		Unix:     unix,
		Name:     name,
		Password: password,
	}
	var resp SigninResponse
	for _, f := range defaultSignin {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	result = &resp
	return
}
