package master

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

var ErrPipeListenerClosed = errors.New(`pipe listener already closed`)

type PipeListener struct {
	ch    chan net.Conn
	close chan struct{}
	done  uint32
	m     sync.Mutex
}

func ListenPipe() *PipeListener {
	return &PipeListener{
		ch:    make(chan net.Conn),
		close: make(chan struct{}),
	}
}

// Accept waits for and returns the next connection to the listener.
func (l *PipeListener) Accept() (c net.Conn, e error) {
	select {
	case c = <-l.ch:
	case <-l.close:
		e = ErrPipeListenerClosed
	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *PipeListener) Close() (e error) {
	if atomic.LoadUint32(&l.done) == 0 {
		l.m.Lock()
		defer l.m.Unlock()
		if l.done == 0 {
			defer atomic.StoreUint32(&l.done, 1)
			close(l.close)
			return
		}
	}
	e = ErrPipeListenerClosed
	return
}

// Addr returns the listener's network address.
func (l *PipeListener) Addr() net.Addr {
	return pipeAddr(0)
}
func (l *PipeListener) Dial(network, addr string) (net.Conn, error) {
	return l.DialContext(context.Background(), network, addr)
}
func (l *PipeListener) DialContext(ctx context.Context, network, addr string) (conn net.Conn, e error) {
	// check closed
	if atomic.LoadUint32(&l.done) != 0 {
		e = ErrPipeListenerClosed
		return
	}

	// pipe
	c0, c1 := net.Pipe()
	// waiting accepted or closed or done
	select {
	case <-ctx.Done():
		e = ctx.Err()
	case l.ch <- c0:
		conn = c1
	case <-l.close:
		c0.Close()
		c1.Close()
		e = ErrPipeListenerClosed
	}
	return
}

type pipeAddr int

func (pipeAddr) Network() string {
	return `pipe`
}
func (pipeAddr) String() string {
	return `pipe`
}
