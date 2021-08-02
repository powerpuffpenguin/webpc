package session

import (
	"context"
	"errors"
)

type PasswordRequest struct {
	Context  context.Context
	ID       int64
	Old      string
	Password string
}
type PasswordResponse struct {
	Changed bool
}

type PasswordHandle func(req *PasswordRequest, resp *PasswordResponse) (e error)

var defaultPassword []PasswordHandle

func ConnectPassword(f PasswordHandle) {
	defaultPassword = append(defaultPassword, f)
}
func Password(ctx context.Context, id int64, old, password string) (changed bool, e error) {
	if len(defaultPassword) == 0 {
		e = errors.New(`signin solt nil`)
		return
	}
	req := &PasswordRequest{
		Context:  ctx,
		ID:       id,
		Old:      old,
		Password: password,
	}
	var resp PasswordResponse
	for _, f := range defaultPassword {
		e = f(req, &resp)
		if e != nil {
			return
		}
	}
	changed = resp.Changed
	return
}
