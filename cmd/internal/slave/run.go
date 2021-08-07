package slave

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/vnet/reverse"
	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	"go.uber.org/zap"
)

type Addr uint8

func (Addr) Network() string {
	return `websocket`
}
func (Addr) String() string {
	return `reverse`
}
func Run(cnf *configure.Connect, system *configure.System, debug bool) {
	var tempDelay time.Duration
	for {
		var d websocket.Dialer
		if cnf.Insecure {
			d.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		l := reverse.Listen(Addr(0), reverse.WithListenerDialContext(func(ctx context.Context, network, address string) (net.Conn, error) {
			c, _, e := d.DialContext(ctx, cnf.URL, nil)
			if e != nil {
				return nil, e
			}
			tempDelay = 0
			return c.UnderlyingConn(), nil
		}))

		srv := newGRPC(&cnf.Option, debug)

		e := srv.Serve(l)
		if e != nil {
			log.Println(e)
		}

		if tempDelay == 0 {
			tempDelay = 5 * time.Second
		} else {
			tempDelay *= 2
		}
		if max := 40 * time.Second; tempDelay > max {
			tempDelay = max
		}
		if ce := logger.Logger.Check(zap.WarnLevel, `Serve error`); ce != nil {
			ce.Write(
				zap.Error(e),
				zap.String(`retrying in`, tempDelay.String()),
			)
		}
		time.Sleep(tempDelay)
	}
}
