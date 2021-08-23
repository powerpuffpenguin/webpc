package vnc

import (
	"context"

	"github.com/powerpuffpenguin/webpc/db"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	ctx, userdata, e := s.JSONUserdata(ctx)
	if e != nil {
		return ctx, e
	} else if userdata.AuthAny(db.Root, db.VNC) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
