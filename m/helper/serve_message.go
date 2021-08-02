package helper

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ServeMessage .
func (Helper) ServeMessage(ctx context.Context, modtime time.Time,
	response func(nobody bool) error,
) (e error) {
	var (
		header      = metadata.MD{}
		method, ims string
	)
	if !isZeroTime(modtime) {
		header.Set(`Last-Modified`, modtime.UTC().Format(http.TimeFormat))
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		strs := md.Get(`Method`)
		if strs != nil && len(strs) != 0 {
			method = strs[0]
		}
		strs = md.Get(`If-Modified-Since`)
		if strs != nil && len(strs) != 0 {
			ims = strs[0]
		}
	}
	if checkIfModifiedSince(method, ims, modtime) == condFalse {
		header.Set(`x-http-code`, strconv.Itoa(http.StatusNotModified))
		e = grpc.SetHeader(ctx, header)
		if e != nil {
			return
		}
		response(true)
		return
	}
	e = grpc.SetHeader(ctx, header)
	if e != nil {
		return
	}
	e = response(method == `HEAD`)
	return
}
