package slave

import (
	"context"
	"fmt"
	"strings"

	"github.com/powerpuffpenguin/webpc/db"

	"google.golang.org/grpc/codes"
)

func (s server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	fmt.Println(fullMethodName)
	ctx, userdata, e := s.Userdata(ctx)
	if e != nil {
		return ctx, e
	} else if userdata.AuthAny(db.Root) {
		return ctx, nil
	}
	if strings.HasSuffix(fullMethodName, `.Find`) {
		return ctx, nil
	}
	return ctx, s.Error(codes.PermissionDenied, `permission denied`)
}
