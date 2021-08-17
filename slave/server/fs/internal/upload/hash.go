package upload

import (
	"encoding/binary"
	"hash/crc32"
	"io"
	"os"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkChunk(chunk uint32) error {
	if chunk < 1024*1024 || chunk > 1024*1024*50 {
		return status.Error(codes.InvalidArgument, `chunk size must rang at [1m,50m]`)
	}
	return nil
}
func checkSize(size int64) error {
	if size < 0 {
		return status.Error(codes.InvalidArgument, `file size not supported 0`)
	}
	return nil
}
func Hash(m *mount.Mount, req *grpc_fs.HashRequest) (resp *grpc_fs.HashResponse, e error) {
	e = checkSize(req.Size)
	if e != nil {
		return
	}
	e = checkChunk(req.Chunk)
	if e != nil {
		return
	}

	f, e := m.Open(req.Path)
	if e != nil {
		return
	}
	hash, match, e := hash(f, req.Size, req.Chunk)
	f.Close()
	if e != nil {
		return
	}
	resp = &grpc_fs.HashResponse{
		Match: match,
		Hash:  hash,
	}
	return
}
func hash(r io.ReadSeeker, size int64, chunk uint32) (val uint32, match bool, e error) {
	offset, e := r.Seek(0, os.SEEK_END)
	if e != nil {
		return
	} else if offset != size {
		return
	}
	match = true
	_, e = r.Seek(0, os.SEEK_CUR)
	if e != nil {
		return
	}
	var (
		hash = crc32.NewIEEE()
		min  = int(chunk)
		b    = make([]byte, chunk)
		n    int
		b32  = make([]byte, 4)
	)

	for {
		n, e = io.ReadAtLeast(r, b, min)
		if n != 0 {
			binary.BigEndian.PutUint32(b32, crc32.ChecksumIEEE(b[:n]))
			hash.Write(b32)
		}
		if e == io.EOF || e == io.ErrUnexpectedEOF {
			break
		} else if e != nil {
			return
		}
	}
	val = hash.Sum32()
	return
}
