package connect

import (
	"io"
	"net"
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	grpc_forward "github.com/powerpuffpenguin/webpc/protocol/forward/forward"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Worker struct {
	server grpc_forward.Forward_ConnectServer
	req    chan *grpc_forward.ConnectRequest
	err    chan error
	close  chan struct{}
}

func New(server grpc_forward.Forward_ConnectServer) *Worker {
	return &Worker{
		server: server,
		req:    make(chan *grpc_forward.ConnectRequest),
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
		} else if req.Event != grpc_forward.Event_Heart {
			select {
			case <-done:
				return
			case w.req <- req:
			}
		}
	}
}
func (w *Worker) chekEvent(evt grpc_forward.Event, expects ...grpc_forward.Event) error {
	for _, expect := range expects {
		if evt == expect {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
}
func (w *Worker) doInit(req *grpc_forward.ConnectRequest) (e error) {
	e = w.chekEvent(req.Event, grpc_forward.Event_Connect)
	if e != nil {
		return
	}
	c, e := net.Dial(`tcp`, req.Addr)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `dial to vnc`); ce != nil {
			ce.Write(
				zap.Error(e),
			)
		}
		return
	}
	e = w.server.Send(&grpc_forward.ConnectResponse{
		Event: grpc_forward.Event_Connect,
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
		e = w.server.Send(&grpc_forward.ConnectResponse{
			Event:  grpc_forward.Event_Binary,
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
		req  *grpc_forward.ConnectRequest
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
			if req.Event == grpc_forward.Event_Binary {
				_, e = c.Write(req.Binary)
				if e != nil {
					work = false
				}
			}
		}
	}
	c.Close()
}
