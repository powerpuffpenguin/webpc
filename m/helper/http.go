package helper

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

var CookieName = strings.ReplaceAll(`github.com/powerpuffpenguin/webpc.session`, `/`, `.`)

func (Helper) CookieName() string {
	return CookieName
}
func (Helper) SetHTTPCookie(ctx context.Context, cookie ...*http.Cookie) error {
	if len(cookie) == 0 {
		return nil
	}
	strs := make([]string, len(cookie))
	for i, c := range cookie {
		strs[i] = c.String()
	}
	md := metadata.MD{}
	md[`set-cookie`] = append(md[`set-cookie`], strs...)
	return grpc.SetHeader(ctx, md)
}
func (Helper) SetHTTPCacheMaxAge(ctx context.Context, maxAge int) error {
	return grpc.SetHeader(ctx, metadata.Pairs(`Cache-Control`, `max-age=`+strconv.Itoa(maxAge)))
}

func (Helper) SetHTTPCode(ctx context.Context, code int) error {
	return grpc.SetHeader(ctx, metadata.Pairs(`x-http-code`, strconv.Itoa(code)))
}

func (h Helper) ToHTTPError(name string, e error) error {
	if os.IsNotExist(e) {
		return h.Error(codes.NotFound, `not exists : `+name)
	}
	if os.IsExist(e) {
		return h.Error(codes.PermissionDenied, `already exists : `+name)
	}
	if os.IsPermission(e) {
		return h.Error(codes.PermissionDenied, `forbidden : `+name)
	}
	return h.Error(codes.Unknown, e.Error())
}
