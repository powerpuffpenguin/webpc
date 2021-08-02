package system

import (
	"context"
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_system "github.com/powerpuffpenguin/webpc/protocol/system"
	"github.com/powerpuffpenguin/webpc/version"
	"time"
)

type server struct {
	grpc_system.UnimplementedSystemServer
	helper.Helper
}

var (
	emptyVersionResponse grpc_system.VersionResponse
	versionResponse      = grpc_system.VersionResponse{
		Platform: version.Platform,
		Version:  version.Version,
	}
)

func (s server) Version(ctx context.Context, req *grpc_system.VersionRequest) (resp *grpc_system.VersionResponse, e error) {
	s.SetHTTPCacheMaxAge(ctx, 60)
	e = s.ServeMessage(ctx,
		time.Unix(startAtResponse.Result, 0),
		func(nobody bool) error {
			if nobody {
				resp = &emptyVersionResponse
			} else {
				resp = &versionResponse
			}
			return nil
		},
	)
	return
}

var (
	emptyStartAtResponse grpc_system.StartAtResponse
	startAtResponse      grpc_system.StartAtResponse
)

func (s server) StartAt(ctx context.Context, req *grpc_system.StartAtRequest) (resp *grpc_system.StartAtResponse, e error) {
	s.SetHTTPCacheMaxAge(ctx, 60)
	e = s.ServeMessage(ctx,
		time.Unix(startAtResponse.Result, 0),
		func(nobody bool) error {
			if nobody {
				resp = &emptyStartAtResponse
			} else {
				resp = &startAtResponse
			}
			return nil
		},
	)
	return
}