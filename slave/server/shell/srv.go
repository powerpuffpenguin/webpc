package shell

import (
	"github.com/powerpuffpenguin/webpc/m/helper"
	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"
)

type server struct {
	grpc_shell.UnimplementedShellServer
	helper.Helper
}
