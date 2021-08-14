package compress

import grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"

type helper struct {
	server grpc_fs.FS_CompressServer
}

func (h helper) SendProgress(name string) error {
	return h.server.Send(&grpc_fs.CompressResponse{
		Event: grpc_fs.Event_Progress,
		Value: name,
	})
}
