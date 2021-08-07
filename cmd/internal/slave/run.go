package slave

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/vnet/reverse"
	"github.com/powerpuffpenguin/webpc/slave"
	"google.golang.org/grpc"
)

type Addr uint8

func (Addr) Network() string {
	return `websocket`
}
func (Addr) String() string {
	return `reverse`
}
func Run() {
	var tempDelay time.Duration
	for {
		var d websocket.Dialer
		l := reverse.Listen(Addr(0), reverse.WithListenerDialContext(func(ctx context.Context, network, address string) (net.Conn, error) {
			c, _, e := d.DialContext(ctx, `ws://127.0.0.1:9000/api/v1/dialer/64048031f73a11eba3890242ac120064`, nil)
			if e != nil {
				return nil, e
			}
			tempDelay = 0
			return c.UnderlyingConn(), nil
		}))

		srv := grpc.NewServer()
		slave.GRPC(srv)

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
		log.Printf("Serve error: %v; retrying in %v", e, tempDelay)
		time.Sleep(tempDelay)
	}
}
