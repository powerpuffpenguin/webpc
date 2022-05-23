package sessionid

import (
	"github.com/powerpuffpenguin/sessionstore/cryptoer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func IsErrExpired(e error) bool {
	if e == nil {
		return false
	}
	if e == cryptoer.ErrExpired {
		return true
	}
	if se, ok := e.(interface {
		GRPCStatus() *status.Status
	}); ok {
		s := se.GRPCStatus()
		if s != nil {
			return s.Code() == codes.Unauthenticated && s.Message() == cryptoer.ErrExpired.Error()
		}
	}
	return false
}
