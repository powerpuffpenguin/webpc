package slave

import (
	"context"
	"log"
	"net"

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
	var d websocket.Dialer
	l := reverse.Listen(Addr(0), reverse.WithListenerDialContext(func(ctx context.Context, network, address string) (net.Conn, error) {
		c, _, e := d.DialContext(ctx, `ws://127.0.0.1:9000/api/v1/dialer`, nil)
		if e != nil {
			log.Println(e)
			return nil, e
		}
		log.Println(`yes`, c.RemoteAddr(), c.LocalAddr())
		return c.UnderlyingConn(), nil
	}))

	srv := grpc.NewServer()
	slave.GRPC(srv)

	e := srv.Serve(l)
	if e != nil {
		log.Fatalln(e)
	}
}
