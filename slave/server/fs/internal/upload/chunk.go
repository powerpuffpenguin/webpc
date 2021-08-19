package upload

import (
	"encoding/hex"
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
		val   string
		items []string
	)
	for i := req.Chunk; i < req.Count; i++ {
		val, e = chunk(m, filepath.Join(dir, strconv.FormatInt(int64(i), 10)))
		if e != nil {
			return
		}
		items = append(items, val)
	}
	resp = &grpc_fs.ChunkResponse{
		Result: items,
	}
	return
}
func chunk(m *mount.Mount, path string) (val string, e error) {
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
	val = hex.EncodeToString(hash.Sum(nil))
	return
}
