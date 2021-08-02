package helper

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Helper int

func (Helper) Error(c codes.Code, msg string) error {
	return status.Error(c, msg)
}
func (Helper) Errorf(c codes.Code, format string, a ...interface{}) error {
	return status.Errorf(c, format, a...)
}
