package open

import (
	"io"
	"os"
	"time"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Worker struct {
	server grpc_fs.FS_OpenServer
	err    chan error
	close  chan struct{}
	file   *os.File
}

func New(server grpc_fs.FS_OpenServer) *Worker {
	return &Worker{
		server: server,
		err:    make(chan error),
		close:  make(chan struct{}),
	}
}
func (w *Worker) mountShared(name string) (m *mount.Mount, e error) {
	fs := mount.Default()
	m = fs.Root(name)
	if m == nil {
		e = status.Error(codes.NotFound, `root not found: `+name)
		return
	}
	if !m.Shared() {
		e = status.Error(codes.PermissionDenied, `filesystem is not shared`)
		return
	}
	return
}
func (w *Worker) recv(evt grpc_fs.FSEvent) (req *grpc_fs.OpenRequest, e error) {
	req, e = w.server.Recv()
	if e != nil {
		return
	}
	if req.Event != evt {
		e = status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
		return
	}
	return
}

type openResult struct {
	File *os.File
	Err  error
}

func (w *Worker) Serve() (e error) {
	ch := make(chan openResult)
	go w.open(ch)

	// recv init event
	t := time.NewTimer(time.Second * 30)
	select {
	case <-t.C:
		e = status.Error(codes.DeadlineExceeded, `wait init timeout`)
	case e = <-w.err:
		if !t.Stop() {
			<-t.C
		}
	case result := <-ch:
		if !t.Stop() {
			<-t.C
		}
		// init
		if result.Err == nil {
			w.file = result.File
			e = w.serve()
			result.File.Close()
		} else {
			e = result.Err
		}
	}

	close(w.close)
	return
}
func (w *Worker) resultOpen(ch chan openResult, f *os.File, e error) {
	select {
	case <-w.close:
		if f != nil {
			f.Close()
		}
	case ch <- openResult{
		File: f,
		Err:  e,
	}:
	}
}
func (w *Worker) open(ch chan openResult) {
	req, e := w.recv(grpc_fs.FSEvent_Open)
	if e != nil {
		w.resultOpen(ch, nil, e)
		return
	}
	// check args
	m, e := w.mountShared(req.Root)
	if e != nil {
		w.resultOpen(ch, nil, e)
		return
	}
	f, e := m.Open(req.Path)
	w.resultOpen(ch, f, e)
}
func (w *Worker) serve() (e error) {
	e = w.server.Send(&grpc_fs.OpenResponse{Event: grpc_fs.FSEvent_Open})
	if e != nil {
		return
	}

	var req grpc_fs.OpenRequest
	for {
		e = w.server.RecvMsg(&req)
		if e != nil {
			if e == io.EOF {
				e = nil
			}
			break
		}
		e = w.onRequest(&req)
		if e != nil {
			break
		}
	}
	return
}
func (w *Worker) onRequest(req *grpc_fs.OpenRequest) (e error) {
	evt := req.Event
	switch evt {
	case grpc_fs.FSEvent_Close:
		e = w.file.Close()
		if e != nil {
			return
		}
		e = w.server.Send(&grpc_fs.OpenResponse{
			Event: evt,
		})
	case grpc_fs.FSEvent_Read:
		e = w.onRead(req.Read)
	case grpc_fs.FSEvent_Seek:
		var ret int64
		ret, e = w.file.Seek(req.Offset, int(req.Whence))
		if e != nil {
			return
		}
		e = w.server.Send(&grpc_fs.OpenResponse{
			Event: evt,
			Seek:  ret,
		})
	case grpc_fs.FSEvent_Readdir:
		var items []os.FileInfo
		items, e = w.file.Readdir(int(req.Readdir))
		if e != nil {
			return
		}
		resp := grpc_fs.OpenResponse{
			Event:   evt,
			Readdir: make([]*grpc_fs.FSInfo, len(items)),
		}
		for i, item := range items {
			resp.Readdir[i] = &grpc_fs.FSInfo{
				Name:    item.Name(),
				Size:    item.Size(),
				Mode:    uint32(item.Mode()),
				ModTime: item.ModTime().Unix(),
				IsDir:   item.IsDir(),
			}
		}
		e = w.server.Send(&resp)
	case grpc_fs.FSEvent_Stat:
		var item os.FileInfo
		item, e = w.file.Stat()
		if e != nil {
			return
		}
		e = w.server.Send(&grpc_fs.OpenResponse{
			Event: evt,
			Stat: &grpc_fs.FSInfo{
				Name:    item.Name(),
				Size:    item.Size(),
				Mode:    uint32(item.Mode()),
				ModTime: item.ModTime().Unix(),
				IsDir:   item.IsDir(),
			},
		})
	default:
		e = status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
	}
	return
}
func (w *Worker) onRead(count uint32) (e error) {
	b := make([]byte, count)
	n, e := w.file.Read(b)
	eof := false
	if e != nil {
		if e == io.EOF {
			eof = true
			e = nil
		} else {
			return
		}
	}
	e = w.server.Send(&grpc_fs.OpenResponse{
		Event: grpc_fs.FSEvent_Read,
		Read:  b[:n],
		Eof:   eof,
	})
	return
}
