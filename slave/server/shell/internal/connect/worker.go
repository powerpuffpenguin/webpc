package connect

import (
	"time"

	"github.com/powerpuffpenguin/webpc/logger"
	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"
	"github.com/powerpuffpenguin/webpc/slave/server/shell/internal/shell"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Worker struct {
	server grpc_shell.Shell_ConnectServer
	req    chan *grpc_shell.ConnectRequest
	err    chan error
	close  chan struct{}

	username string
}

func New(server grpc_shell.Shell_ConnectServer, username string) *Worker {
	return &Worker{
		server:   server,
		username: username,
		req:      make(chan *grpc_shell.ConnectRequest),
		err:      make(chan error),
		close:    make(chan struct{}),
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
		} else if req.Event != grpc_shell.Event_Heart {
			select {
			case <-done:
				return
			case w.req <- req:
			}
		}
	}
}
func (w *Worker) chekEvent(evt grpc_shell.Event, expects ...grpc_shell.Event) error {
	for _, expect := range expects {
		if evt == expect {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
}
func (w *Worker) doInit(req *grpc_shell.ConnectRequest) (e error) {
	e = w.chekEvent(req.Event, grpc_shell.Event_Connect)
	if e != nil {
		return
	}
	var (
		newshell bool
		shellid  int64
		conn     = newConn(w.server)
	)
	if req.Id == 0 {
		newshell = true
		shellid = time.Now().UTC().Unix()
	} else {
		shellid = req.Id
	}
	s, e := shell.DefaultManager().Attach(conn, w.username, shellid, uint16(req.Cols), uint16(req.Rows), newshell)
	if e != nil {
		return
	}
	w.serve(s, conn)
	s.Unattack(conn)
	conn.Close()
	return
}
func (w *Worker) serve(s *shell.Shell, c *conn) (e error) {
	var (
		req  *grpc_shell.ConnectRequest
		resp *grpc_shell.ConnectResponse
	)
	for {
		select {
		case e = <-w.err:
			return
		case req = <-w.req:
			w.onRequest(s, req)
		case resp = <-c.resp:
			e = w.server.Send(resp)
			if e != nil {
				return
			}
		case <-c.close:
			return
		}
	}
}
func (w *Worker) onRequest(s *shell.Shell, req *grpc_shell.ConnectRequest) {
	var (
		err error
		evt = req.Event
	)
	switch evt {
	case grpc_shell.Event_Binary:
		_, err = s.Write(req.Binary)
	case grpc_shell.Event_Resize:
		err = s.SetSize(uint16(req.Cols), uint16(req.Rows))
	case grpc_shell.Event_FontSize:
		err = s.SetFontSize(req.FontSize)
	case grpc_shell.Event_FontFamily:
		err = s.SetFontFamily(req.FontFamily)
	default:
		err = status.Error(codes.InvalidArgument, `unexpected event: `+evt.String())
	}
	if err != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `shell event error`); ce != nil {
			ce.Write(
				zap.Error(err),
				zap.String(`event`, req.Event.String()),
			)
		}
	}
}
