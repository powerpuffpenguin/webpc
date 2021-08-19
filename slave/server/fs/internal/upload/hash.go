package upload

import (
	"encoding/binary"
	"encoding/hex"
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

var emptyHashResponse grpc_fs.HashResponse

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
		if status.Code(e) == codes.NotFound {
			e = nil
			resp = &emptyHashResponse
		}
		return
	}
	hash, e := hash(f, req.Size, req.Chunk)
	f.Close()
	if e != nil {
		return
	}
	resp = &grpc_fs.HashResponse{
		Exists: true,
		Hash:   hash,
	}
	return
}
func hash(r io.ReadSeeker, size int64, chunk uint32) (val string, e error) {
	offset, e := r.Seek(0, os.SEEK_END)
	if e != nil {
		return
	} else if offset != size {
		return
	}
	_, e = r.Seek(0, os.SEEK_SET)
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
			e = nil
			break
		} else if e != nil {
			return
		}
	}
	val = hex.EncodeToString(hash.Sum(nil))
	return
}
