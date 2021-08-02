package fs

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/powerpuffpenguin/webpc/m/helper"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
)

var name string

func init() {
	rand.Seed(time.Now().Unix())
	name = fmt.Sprint(rand.Int())
}

type server struct {
	grpc_fs.UnimplementedFSServer
	helper.Helper
}

func (server) Mount(ctx context.Context, req *grpc_fs.MountRequest) (resp *grpc_fs.MountResponse, e error) {
	resp = &grpc_fs.MountResponse{
		Name: []string{
			name,
		},
	}
	return
}
