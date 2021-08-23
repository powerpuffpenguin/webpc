package connect

import (
	"io"
	"net"
	"time"

	grpc_vnc "github.com/powerpuffpenguin/webpc/protocol/forward/vnc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var VNC = ``

type Worker struct {
	server grpc_vnc.Vnc_ConnectServer
	req    chan *grpc_vnc.ConnectRequest
	err    chan error
	close  chan struct{}
}

func New(server grpc_vnc.Vnc_ConnectServer) *Worker {
	return &Worker{
		server: server,
		req:    make(chan *grpc_vnc.ConnectRequest),
		err:    make(chan error),
		close:  make(chan struct{}),
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
		} else if req.Event != grpc_vnc.Event_Heart {
			select {
			case <-done:
				return
			case w.req <- req:
			}
		}
	}
}
func (w *Worker) chekEvent(evt grpc_vnc.Event, expects ...grpc_vnc.Event) error {
	for _, expect := range expects {
		if evt == expect {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
}
func (w *Worker) doInit(req *grpc_vnc.ConnectRequest) (e error) {
	e = w.chekEvent(req.Event, grpc_vnc.Event_Connect)
	if e != nil {
		return
	}
	c, e := net.Dial(`tcp`, VNC)
	if e != nil {
		return
	}
	e = w.server.Send(&grpc_vnc.ConnectResponse{
		Event: grpc_vnc.Event_Connect,
	})
	if e != nil {
		c.Close()
		return
	}
	e = w.serve(c)
	return
}
func (w *Worker) serve(c net.Conn) (e error) {
	go w.toVNC(c)
	var (
		b = make([]byte, 1024*32)
		n int
	)
	for {
		n, e = c.Read(b)
		if e != nil {
			if e == io.EOF {
				e = nil
			}
			break
		}
		e = w.server.Send(&grpc_vnc.ConnectResponse{
			Event:  grpc_vnc.Event_Binary,
			Binary: b[:n],
		})
		if e != nil {
			break
		}
	}
	c.Close()
	return
}
func (w *Worker) toVNC(c net.Conn) {
	var (
		req  *grpc_vnc.ConnectRequest
		e    error
		work = true
	)
	for work {
		select {
		case <-w.err:
			work = false
		case <-w.close:
			work = false
		case req = <-w.req:
			if req.Event == grpc_vnc.Event_Binary {
				_, e = c.Write(req.Binary)
				if e != nil {
					work = false
				}
			}
		}
	}
	c.Close()
}
