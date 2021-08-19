package upload

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strconv"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var emptyMergeResponse grpc_fs.MergeResponse

func Merge(m *mount.Mount, req *grpc_fs.MergeRequest) (resp *grpc_fs.MergeResponse, e error) {
	dir, name := filepath.Split(req.Path)
	dir = filepath.Join(dir, `.chunks_`+name)
	dst := filepath.Join(dir, `ok`)

	f, e := m.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if e != nil {
		return
	}
	e = merge(m, dir, req.Hash, uint64(req.Count), f)
	f.Close()
	if e != nil {
		return
	}
	// rename
	oldpath, e := m.Filename(dst)
	if e != nil {
		return
	}
	newpath, e := m.Filename(req.Path)
	if e != nil {
		return
	}
	e = os.Rename(oldpath, newpath)
	if e != nil {
		if os.IsNotExist(e) {
			e = status.Error(codes.NotFound, `not exists: `+dst)
		} else if os.IsExist(e) {
			e = status.Error(codes.AlreadyExists, `already exists: `+req.Path)
		} else if os.IsPermission(e) {
			e = status.Error(codes.PermissionDenied, `forbidden move: `+dst+` -> `+req.Path)
		}
		return
	}
	resp = &emptyMergeResponse
	m.RemoveAll(dir)
	return
}
func merge(m *mount.Mount, dir string, val uint32, count uint64, w io.Writer) (e error) {
	var (
		hash  = crc32.NewIEEE()
		chunk = crc32.NewIEEE()
		b32   = make([]byte, 4)
		r     io.ReadCloser
	)

	for i := uint64(0); i < count; i++ {
		chunk.Reset()

		name := filepath.Join(dir, strconv.FormatUint(i, 10))
		r, e = m.Open(name)
		if e != nil {
			return
		}
		_, e = io.Copy(io.MultiWriter(chunk, w), r)
		r.Close()
		if e != nil {
			return
		}
		binary.BigEndian.PutUint32(b32, chunk.Sum32())
		hash.Write(b32)
	}
	if hash.Sum32() != val {
		e = status.Error(codes.InvalidArgument, `hash not match`)
		return
	}
	return
}
