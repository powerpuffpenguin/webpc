package dialer

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/powerpuffpenguin/vnet"
)

type Listener struct {
	ch    chan *Conn
	close chan struct{}
	done  uint32
	m     sync.Mutex
}

func newListener() *Listener {
	return &Listener{}
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	select {
	case c := <-l.ch:
		return c, nil
	case <-l.close:
		return nil, vnet.ErrListenerClosed
	}
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() (e error) {
	if atomic.LoadUint32(&l.done) == 0 {
		l.m.Lock()
		defer l.m.Unlock()
		if l.done == 0 {
			defer atomic.StoreUint32(&l.done, 1)
			close(l.ch)
			return
		}
	}
	e = vnet.ErrListenerClosed
	return
}

// Addr returns the listener's network address.
func (*Listener) Addr() net.Addr {
	return addr(0)
}

type addr uint

func (addr) Network() string {
	return `websocket`
}
func (addr) String() string {
	return `dialer`
}
