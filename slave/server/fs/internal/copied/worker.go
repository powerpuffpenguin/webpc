package copied

import (
	"time"

	grpc_fs "github.com/powerpuffpenguin/webpc/protocol/forward/fs"
	"github.com/powerpuffpenguin/webpc/single/mount"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Worker struct {
	server   grpc_fs.FS_CopyServer
	req      chan *grpc_fs.CopyRequest
	readable bool
	err      chan error
	close    chan struct{}
	style    grpc_fs.Event

	SrcRoot string
	SrcDir  string
	DstRoot string
	DstDir  string
	Names   []string
	Copied  bool
}

func New(server grpc_fs.FS_CopyServer, readable bool) *Worker {
	return &Worker{
		server:   server,
		readable: readable,
		req:      make(chan *grpc_fs.CopyRequest),
		err:      make(chan error),
		close:    make(chan struct{}),
		style:    grpc_fs.Event_EventUniversal,
	}
}

func (w *Worker) mountWrite(name string) (m *mount.Mount, e error) {
	fs := mount.Default()
	m = fs.Root(name)
	if m == nil {
		e = status.Error(codes.NotFound, `root not found: `+name)
		return
	}
	if !m.Write() {
		e = status.Error(codes.PermissionDenied, `filesystem is not writable`)
		return
	}
	return
}
func (w *Worker) mountRead(name string) (m *mount.Mount, e error) {
	fs := mount.Default()
	m = fs.Root(name)
	if m == nil {
		e = status.Error(codes.NotFound, `root not found: `+name)
		return
	}
	if w.readable {
		if !m.Read() {
			e = status.Error(codes.PermissionDenied, `filesystem is not readable`)
			return
		}
	} else {
		if !m.Shared() {
			e = status.Error(codes.PermissionDenied, `user forbidden to read`)
			return
		}
	}
	return
}
func (w *Worker) chekEvent(evt grpc_fs.Event, expects ...grpc_fs.Event) error {
	for _, expect := range expects {
		if evt == expect {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
}
func (w *Worker) waitRequest(expects ...grpc_fs.Event) (req *grpc_fs.CopyRequest, e error) {
	select {
	case req = <-w.req:
	case e = <-w.err:
		return
	}
	e = w.chekEvent(req.Event, expects...)
	if e != nil {
		req = nil
		return
	}
	return
}
func (w *Worker) runRecv() {
	done := w.close
	server := w.server
	for {
		req, e := server.Recv()
		if e != nil {
			select {
			case w.err <- e:
			case <-done:
			}
			break
		} else if req.Event != grpc_fs.Event_Heart {
			select {
			case <-done:
				return
			case w.req <- req:
			}
		}
	}
}
func (w *Worker) Serve() (e error) {
	go w.runRecv()

	// recv init event
	t := time.NewTimer(time.Second * 30)
	select {
	case <-t.C:
		e = status.Error(codes.DeadlineExceeded, `wait init timeout`)
	case e = <-w.err:
		if !t.Stop() {
			<-t.C
		}
	case req := <-w.req:
		if !t.Stop() {
			<-t.C
		}
		// init
		e = w.doInit(req)
	}

	close(w.close)
	return
}
func (w *Worker) doInit(req *grpc_fs.CopyRequest) (e error) {
	e = w.chekEvent(req.Event, grpc_fs.Event_Init)
	if e != nil {
		return
	}
	// check args
	mw, e := w.mountWrite(req.SrcRoot)
	if e != nil {
		return
	}
	_, e = mw.Filename(req.SrcDir)
	if e != nil {
		return
	}
	mr, e := w.mountRead(req.DstRoot)
	if e != nil {
		return
	}
	_, e = mr.Filename(req.SrcDir)
	if e != nil {
		return
	}
	for _, name := range req.Names {
		e = mw.CheckName(name)
		if e != nil {
			return
		}
	}

	w.SrcRoot = req.SrcRoot
	w.SrcDir = req.SrcDir
	w.DstRoot = req.DstRoot
	w.DstDir = req.DstDir
	w.Names = req.Names
	w.Copied = req.Copied
	e = w.serve(mw, mr)
	if e == nil {
		w.server.Send(&grpc_fs.CopyResponse{
			Event: grpc_fs.Event_Success,
		})
	}
	return
}
