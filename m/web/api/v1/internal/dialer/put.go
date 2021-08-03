package dialer

import (
	"context"
	"net"

	"sync"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/vnet/reverse"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/forward"
	"github.com/powerpuffpenguin/webpc/slave"
	"github.com/powerpuffpenguin/webpc/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func Put(c *Conn) {
	defaultListeners.ch <- c
}

var defaultListeners = newListeners()

type Element struct {
	dialer   *reverse.Dialer
	listener *Listener
	id       int64
	gateway  *runtime.ServeMux
	cc       *grpc.ClientConn
}

type Listeners struct {
	ch   chan *Conn
	keys map[int64]*Element
	sync.Mutex
}

func newListeners() *Listeners {
	l := &Listeners{
		ch:   make(chan *Conn),
		keys: make(map[int64]*Element),
	}
	go l.work()
	return l
}
func (l *Listeners) work() {
	for {
		l.onConn(<-l.ch)
	}
}
func (l *Listeners) onConn(c *Conn) {
	ele, exists := l.keys[c.id]
	if !exists {
		listener := newListener()
		ele = &Element{
			id:       c.id,
			listener: listener,
			dialer:   reverse.NewDialer(listener),
		}
		go ele.dialer.Serve()
		cc, e := grpc.Dial(``,
			grpc.WithInsecure(),
			grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
				return ele.dialer.DialContext(c, `reverse`, s)
			}))
		if e != nil {
			if ce := logger.Logger.Check(zap.ErrorLevel, `grpc.Dial error`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			c.Close()
			return
		}
		gateway := utils.NewGateway()
		ele.gateway = gateway
		ele.cc = cc
		e = slave.HTTP(gateway, cc)
		if e != nil {
			if ce := logger.Logger.Check(zap.ErrorLevel, `register gateway error`); ce != nil {
				ce.Write(
					zap.Error(e),
				)
			}
			c.Close()
			cc.Close()
			return
		}
		forward.Default().Put(c.id, cc, gateway)
		l.keys[c.id] = ele
	}
	go l.putConn(ele, c)
}
func (l *Listeners) putConn(ele *Element, c *Conn) {
	e := ele.listener.Put(c)
	if e == nil {
		return
	}
	c.Close()
}
