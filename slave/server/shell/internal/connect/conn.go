package connect

import (
	"errors"
	"sync/atomic"

	grpc_shell "github.com/powerpuffpenguin/webpc/protocol/forward/shell"
)

var errAlreadyClosed = errors.New(`conn already closed`)

type conn struct {
	server grpc_shell.Shell_ConnectServer
	close  chan struct{}
	closed int32
	resp   chan *grpc_shell.ConnectResponse
}

func newConn(server grpc_shell.Shell_ConnectServer) *conn {
	return &conn{
		server: server,
		close:  make(chan struct{}),
		resp:   make(chan *grpc_shell.ConnectResponse),
	}
}
func (c *conn) Close() (e error) {
	if c.closed == 0 && atomic.SwapInt32(&c.closed, 1) == 0 {
		close(c.close)
	} else {
		e = errAlreadyClosed
	}
	return
}
func (c *conn) WriteBinary(b []byte) (e error) {
	select {
	case <-c.close:
		e = errAlreadyClosed
	case c.resp <- &grpc_shell.ConnectResponse{
		Event:  grpc_shell.Event_Binary,
		Binary: b,
	}:
	}
	return
}
func (c *conn) WriteString(str string) (e error) {
	select {
	case <-c.close:
		e = errAlreadyClosed
	case c.resp <- &grpc_shell.ConnectResponse{
		Event:  grpc_shell.Event_Binary,
		Binary: []byte(str),
	}:
	}
	return
}
func (c *conn) WriteInfo(id int64, name string, started int64, fontSize int32, fontFamily string) (e error) {
	select {
	case <-c.close:
		e = errAlreadyClosed
	case c.resp <- &grpc_shell.ConnectResponse{
		Event:      grpc_shell.Event_Info,
		Id:         id,
		Name:       name,
		At:         started,
		FontSize:   fontSize,
		FontFamily: fontFamily,
	}:
	}
	return
}
