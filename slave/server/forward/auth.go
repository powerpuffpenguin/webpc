package forward

import (
	"context"

	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/db"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if !configure.DefaultSystem().PortForward {
		return ctx, s.Error(codes.PermissionDenied, `port forward not enable`)
	}

	ctx, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return ctx, e
	} else if userdata.AuthAny(db.Root, db.PortForward) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
