package forward

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/powerpuffpenguin/webpc/cmd/internal/client"
	socks5_impl "github.com/powerpuffpenguin/webpc/cmd/internal/forward/socks5"
	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
)

type Worker struct {
	dialer *client.Dialer
}

func newWorker(dialer *client.Dialer) *Worker {
	return &Worker{
		dialer: dialer,
	}
}
func (w *Worker) Serve(l net.Listener, remote string, socks5 bool) {
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
				if ce := logger.Logger.Check(zap.WarnLevel, `Accept error`); ce != nil {
					ce.Write(
						zap.Error(e),
						zap.Duration(`retrying in`, tempDelay),
					)
				}
				time.Sleep(tempDelay)
				continue
			}
			break
		}
		tempDelay = 0
		go w.serve(c, remote, socks5)
	}
	l.Close()
}
func (w *Worker) serve(c0 net.Conn, remote string, socks5 bool) {
	var fromAddr string
	if ce := logger.Logger.Check(zap.DebugLevel, `one in`); ce != nil {
		fromAddr = c0.RemoteAddr().String()
		ce.Write(
			zap.String(`from`, fromAddr),
			zap.String(`to`, remote),
		)
		at := time.Now()
		defer func() {
			logger.Logger.Debug(`one out`,
				zap.String(`from`, fromAddr),
				zap.String(`to`, remote),
				zap.Duration(`duration`, time.Since(at)),
			)
		}()
	}

	var (
		ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
		e           error
		c1          net.Conn
		version     byte
	)
	defer cancel()
	if socks5 {
		version, remote, e = socks5_impl.Recv(ctx, c0)
		if e != nil {
			c0.Close()
			return
		}
	}
	c1, e = w.dialer.DialContext(ctx, `tcp`, remote)
	if e != nil {
		if ce := logger.Logger.Check(zap.WarnLevel, `dial error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`from`, fromAddr),
				zap.String(`to`, remote),
			)
		}
		c0.Close()
		return
	}
	if ce := logger.Logger.Check(zap.DebugLevel, `dial success`); ce != nil {
		ce.Write(
			zap.String(`from`, fromAddr),
			zap.String(`to`, remote),
		)
	}

	if socks5 {
		e = socks5_impl.Send(c0, version)
		if e != nil {
			c0.Close()
			c1.Close()
			return
		}
	}
	go w.forward(c1, c0)
	w.forward(c0, c1)
}
func (w *Worker) forward(c0, c1 net.Conn) {
	io.Copy(c0, c1)
	c0.Close()
	c1.Close()
}
