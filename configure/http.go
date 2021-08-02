package configure

import (
	"time"

	"google.golang.org/grpc/keepalive"
)

// HTTP configure
type HTTP struct {
	Addr string

	CertFile string
	KeyFile  string

	Swagger bool

	Option ServerOption
}

// H2 if tls return true
func (h *HTTP) H2() bool {
	return h.CertFile != `` && h.KeyFile != ``
}

// H2C if not use tls return true
func (h *HTTP) H2C() bool {
	return h.CertFile == `` || h.KeyFile == ``
}

type ServerOption struct {
	WriteBufferSize, ReadBufferSize          int
	InitialWindowSize, InitialConnWindowSize int32
	MaxRecvMsgSize, MaxSendMsgSize           int
	MaxConcurrentStreams                     uint32
	ConnectionTimeout                        time.Duration
	Keepalive                                keepalive.ServerParameters
}
