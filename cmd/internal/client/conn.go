package client

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	grpc_forward "github.com/powerpuffpenguin/webpc/protocol/forward/forward"
)

type Conn struct {
	*websocket.Conn
	ticker *time.Ticker
	close  chan struct{}
	buf    []byte
	m      sync.Mutex
}

func NewConn(ws *websocket.Conn, heart int) (conn *Conn, e error) {
	var ticker *time.Ticker
	close := make(chan struct{})
	if heart > 0 {
		var b []byte
		b, e = Marshal(&grpc_forward.ConnectRequest{
			Event: grpc_forward.Event_Heart,
		})
		if e != nil {
			return
		}
		ticker = time.NewTicker(time.Second * time.Duration(heart))
		go func() {
			for {
				select {
				case <-close:
					return
				case <-ticker.C:
					ws.WriteMessage(websocket.TextMessage, b)
				}
			}
		}()
	}
	conn = &Conn{
		Conn:   ws,
		ticker: ticker,
		close:  close,
	}
	return
}
func (c *Conn) Read(b []byte) (n int, e error) {
	var t int
	c.m.Lock()
	buf := c.buf
	if len(buf) == 0 {
		for {
			t, buf, e = c.ReadMessage()
			if e != nil {
				return
			} else if t == websocket.BinaryMessage &&
				len(b) != 0 {
				break
			}
		}
	}
	n = copy(b, buf)
	buf = buf[n:]

	if len(buf) == 0 {
		c.buf = nil
	} else {
		c.buf = buf
	}
	c.m.Unlock()
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if err == nil {
		n = len(b)
	}
	return
}
func (c *Conn) Close() error {
	e := c.Conn.Close()
	if e == nil {
		if c.ticker != nil {
			c.ticker.Stop()
		}
		close(c.close)
	}
	return e
}
func (c *Conn) SetDeadline(t time.Time) error {
	return c.UnderlyingConn().SetDeadline(t)
}
