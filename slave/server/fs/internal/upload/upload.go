package upload

import (
	"os"
	"path/filepath"
	"strconv"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
)

var emptyUploadResponse grpc_fs.UploadResponse

func Upload(m *mount.Mount, req *grpc_fs.UploadRequest) (resp *grpc_fs.UploadResponse, e error) {
	dir, name := filepath.Split(req.Path)
	dir = filepath.Join(dir, `.chunks_`+name)
	os.Mkdir(dir, 0775)
	path := filepath.Join(dir, strconv.FormatInt(int64(req.Chunk), 10))
	f, e := m.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		return
	}
	_, e = f.Write(req.Data)
	f.Close()
	if e != nil {
		return
	}
	resp = &emptyUploadResponse
	return
}
