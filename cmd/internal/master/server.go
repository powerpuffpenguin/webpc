package master

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/powerpuffpenguin/webpc/configure"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/register"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type Server struct {
	pipe  *PipeListener
	gpipe *grpc.Server

	tcp  net.Listener
	gtcp *grpc.Server

	mux *gin.Engine
}

func newServer(l net.Listener, swagger, debug bool, cnf *configure.ServerOption) (s *Server) {
	pipe := ListenPipe()
	clientConn, e := grpc.Dial(`pipe`,
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(c context.Context, s string) (net.Conn, error) {
			return pipe.DialContext(c, `pipe`, s)
		}),
	)
	if e != nil {
		logger.Logger.Panic(`pipe`,
			zap.Error(e),
		)
	}

	gateway := newGateway()
	mux := gin.Default()
	register.HTTP(clientConn, mux, gateway, swagger)

	s = &Server{
		pipe:  pipe,
		gpipe: newGRPC(cnf, gateway, clientConn, debug),
		tcp:   l,
		gtcp:  newGRPC(cnf, nil, nil, debug),
		mux:   mux,
	}
	return
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextType := r.Header.Get(`Content-Type`)
	if r.ProtoMajor == 2 && strings.Contains(contextType, `application/grpc`) {
		s.gtcp.ServeHTTP(w, r)
	} else {
		s.mux.ServeHTTP(w, r)
	}
}
func (s *Server) Serve() (e error) {
	go s.gpipe.Serve(s.pipe)

	// h2c
	var httpServer http.Server
	var http2Server http2.Server
	e = http2.ConfigureServer(&httpServer, &http2Server)
	if e != nil {
		return
	}
	httpServer.Handler = h2c.NewHandler(s, &http2Server)
	// h2c Serve
	e = httpServer.Serve(s.tcp)
	return
}
func (s *Server) ServeTLS(certFile, keyFile string) (e error) {
	go s.gpipe.Serve(s.pipe)

	e = http.ServeTLS(s.tcp, s, certFile, keyFile)
	return
}
