package system

import (
	"context"
	"time"

	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_system "github.com/powerpuffpenguin/webpc/protocol/forward/system"
	"github.com/powerpuffpenguin/webpc/single/upgrade"
	"github.com/powerpuffpenguin/webpc/version"
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
		Commit:   version.Commit,
		Date:     version.Date,
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
	startAtResponse      = grpc_system.StartAtResponse{
		Result: time.Now().Unix(),
	}
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

var emptyUpgradedResponse grpc_system.UpgradedResponse

func (s server) Upgraded(ctx context.Context, req *grpc_system.UpgradedRequest) (resp *grpc_system.UpgradedResponse, e error) {
	modtime, version, upgraded := upgrade.DefaultUpgrade().Upgraded()
	s.SetHTTPCacheMaxAge(ctx, 60*60)
	e = s.ServeMessage(ctx,
		modtime,
		func(nobody bool) error {
			if nobody || !upgraded {
				resp = &emptyUpgradedResponse
			} else {
				resp = &grpc_system.UpgradedResponse{
					Version: version,
				}
			}
			return nil
		},
	)
	return
}
