package dialer

import (
	"sync"

	"github.com/powerpuffpenguin/vnet/reverse"
)

func Put(c *Conn) {
	defaultListeners.ch <- c
}

var defaultListeners = newListeners()

type Element struct {
	listener *Listener
	id       int64
	num      int
}

type Listeners struct {
	ch   chan *Conn
	keys map[int64]*Element
	sync.Mutex
}

func newListeners() *Listeners {
	l := &Listeners{
		ch: make(chan *Conn),
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
	// f := register.DefaultForward()
	ele, exists := l.keys[c.id]
	if exists {
		ele.num++
	} else {
		ele = &Element{
			id:       c.id,
			num:      1,
			listener: newListener(),
		}
		reverse.NewDialer(ele.listener)

		l.keys[c.id] = ele
	}
}
