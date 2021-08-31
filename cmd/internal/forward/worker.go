package forward

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/powerpuffpenguin/webpc/cmd/internal/client"
)

type Worker struct {
	dialer *client.Dialer
}

func newWorker(dialer *client.Dialer) *Worker {
	return &Worker{
		dialer: dialer,
	}
}
func (w *Worker) Serve(l net.Listener, remote string) {
	var tempDelay time.Duration
	for {
		c, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			break
		}
		tempDelay = 0
		go w.serve(c, remote)
	}
	l.Close()
}
func (w *Worker) serve(c0 net.Conn, remote string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	c1, e := w.dialer.DialContext(ctx, `tcp`, remote)
	if e != nil {
		c0.Close()
		return
	}
	go w.forward(c1, c0)
	w.forward(c0, c1)
}
func (w *Worker) forward(c0, c1 net.Conn) {
	io.Copy(c0, c1)
	c0.Close()
	c1.Close()
}
