package dialer

import (
	"net"
)

type Conn struct {
	id int64
	net.Conn
	close    chan struct{}
	callback func(id int64)
}

func NewConn(id int64, c net.Conn) *Conn {
	return &Conn{
		id:    id,
		Conn:  c,
		close: make(chan struct{}),
	}
}
func (c *Conn) Close() (e error) {
	e = c.Conn.Close()
	if e == nil {
		close(c.close)
		if c.callback != nil {
			c.callback(c.id)
		}
	}
	return
}
func (c *Conn) Done() <-chan struct{} {
	return c.close
}
