package dialer

import (
	"context"
	"net"
	"time"

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

const (
	CheckTicker   = time.Minute
	CheckDeadline = CheckTicker * 8
)

func Put(c *Conn) {
	defaultListeners.ch <- c
}

var defaultListeners = newListeners()

type Element struct {
	dialer         *reverse.Dialer
	listener       *Listener
	cc             *grpc.ClientConn
	gateway        *runtime.ServeMux
	referenceCount uint64
	deadline       time.Time
}

// Listeners gateway splitter
type Listeners struct {
	ch   chan *Conn
	exit chan int64

	// ready listener
	ready map[int64]*Element
	// idle listener
	idle map[int64]*Element
	sync.Mutex
}

func newListeners() *Listeners {
	l := &Listeners{
		ch:    make(chan *Conn),
		exit:  make(chan int64),
		ready: make(map[int64]*Element),
		idle:  make(map[int64]*Element),
	}
	go l.work()
	return l
}
func (l *Listeners) work() {
	t := time.NewTicker(CheckTicker)
	for {
		select {
		case c := <-l.ch:
			l.onConn(c)
		case id := <-l.exit:
			l.onExit(id)
		case at := <-t.C:
			l.onCheck(at)
		}
	}
}
func (l *Listeners) onCheck(at time.Time) {
	for key, ele := range l.idle {
		if at.Before(ele.deadline) {
			continue
		}
		ele.dialer.Close()
		ele.cc.Close()
		delete(l.idle, key)
		if ce := logger.Logger.Check(zap.InfoLevel, `delete listener`); ce != nil {
			ce.Write(
				zap.Int64(`id`, key),
			)
		}
	}
}
func (l *Listeners) onExit(id int64) {
	ele := l.ready[id]
	ele.referenceCount--
	if ele.referenceCount == 0 {
		forward.Default().Del(id)

		delete(l.ready, id)
		l.idle[id] = ele
		ele.deadline = time.Now().Add(CheckDeadline)
		if ce := logger.Logger.Check(zap.InfoLevel, `set listener idle`); ce != nil {
			ce.Write(
				zap.Int64(`id`, id),
			)
		}
	}
}
func (l *Listeners) getElement(c *Conn) (ele *Element, e error) {
	var (
		id     = c.id
		exists bool
	)
	// find from idle
	ele, exists = l.idle[id]
	if exists {
		if ce := logger.Logger.Check(zap.InfoLevel, `set listener ready`); ce != nil {
			ce.Write(
				zap.Int64(`id`, id),
				zap.String(`remote`, c.RemoteAddr().String()),
			)
		}
		delete(l.idle, id)
		l.ready[id] = ele
		ele.referenceCount = 1
		forward.Default().Put(id, ele.cc, ele.gateway)
		return
	}
	// find from ready
	ele, exists = l.ready[id]
	if exists {
		ele.referenceCount++
		if ce := logger.Logger.Check(zap.DebugLevel, `listener++`); ce != nil {
			ce.Write(
				zap.Int64(`id`, id),
				zap.String(`remote`, c.RemoteAddr().String()),
				zap.Uint64(`reference count`, ele.referenceCount),
			)
		}
		return
	}
	// create new
	ele, e = l.createListener(id)
	if e != nil {
		if ce := logger.Logger.Check(zap.InfoLevel, `create listener error`); ce != nil {
			ce.Write(
				zap.Int64(`id`, id),
				zap.String(`remote`, c.RemoteAddr().String()),
			)
		}
		return
	}
	l.ready[id] = ele
	forward.Default().Put(id, ele.cc, ele.gateway)
	if ce := logger.Logger.Check(zap.InfoLevel, `create listener`); ce != nil {
		ce.Write(
			zap.Int64(`id`, id),
			zap.String(`remote`, c.RemoteAddr().String()),
		)
	}
	return
}
func (l *Listeners) createListener(id int64) (ele *Element, e error) {
	listener := newListener()
	dialer := reverse.NewDialer(listener)
	// grpc client
	cc, e := grpc.Dial(``,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return dialer.DialContext(c, `reverse`, s)
		}),
	)
	if e != nil {
		dialer.Close()
		return
	}
	// register gateway
	gateway := utils.NewGateway()
	e = slave.HTTP(gateway, cc)
	if e != nil {
		dialer.Close()
		cc.Close()
		return
	}

	// create
	ele = &Element{
		listener:       listener,
		dialer:         dialer,
		cc:             cc,
		gateway:        gateway,
		referenceCount: 1,
	}
	go dialer.Serve()
	return
}
func (l *Listeners) onConn(c *Conn) {
	ele, e := l.getElement(c)
	if e != nil {
		c.Close()
		return
	}
	c.callback = l.onCallback
	go l.putConn(ele, c)
}
func (l *Listeners) onCallback(id int64) {
	l.exit <- id
}
func (l *Listeners) putConn(ele *Element, c *Conn) {
	e := ele.listener.Put(c)
	if e == nil {
		return
	}
	c.Close()
}
