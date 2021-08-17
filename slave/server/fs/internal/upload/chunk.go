package upload

import (
	"hash/crc32"
	"io"
	"path/filepath"
	"strconv"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Chunk(m *mount.Mount, req *grpc_fs.ChunkRequest) (resp *grpc_fs.ChunkResponse, e error) {
	dir, name := filepath.Split(req.Path)
	dir = filepath.Join(dir, `.chunks_`+name)
	var (
		val    uint32
		exists bool
		items  []*grpc_fs.ChunkData
	)
	for _, index := range req.Chunk {
		val, exists, e = chunk(m, filepath.Join(dir, strconv.FormatInt(int64(index), 10)))
		if e != nil {
			return
		}
		items = append(items, &grpc_fs.ChunkData{
			Exists: exists,
			Hash:   val,
		})
	}
	resp = &grpc_fs.ChunkResponse{
		Result: items,
	}
	return
}
func chunk(m *mount.Mount, path string) (val uint32, exists bool, e error) {
	f, e := m.Open(path)
	if e != nil {
		if status.Code(e) == codes.NotFound {
			e = nil
		}
		return
	}
	hash := crc32.NewIEEE()
	_, e = io.Copy(hash, f)
	f.Close()
	if e != nil {
		return
	}
	exists = true
	val = hash.Sum32()
	return
}
