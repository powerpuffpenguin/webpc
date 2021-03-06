package vnc

import (
	"context"

	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/slave/server/vnc/internal/connect"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if connect.VNC == `` {
		return ctx, s.Error(codes.PermissionDenied, `vnc not enable`)
	}

	ctx, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return ctx, e
	} else if userdata.AuthAny(db.Root, db.VNC) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
