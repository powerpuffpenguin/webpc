package slave

import (
	"context"
	"strings"

	"github.com/powerpuffpenguin/webpc/db"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	ctx, userdata, e := s.Userdata(ctx)
	if e != nil {
		return ctx, e
	} else if strings.HasSuffix(fullMethodName, `/Find`) ||
		userdata.AuthAny(db.Root) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
