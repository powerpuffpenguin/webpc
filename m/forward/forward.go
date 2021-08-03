package forward

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/m/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var defaultForward = newForward()

func Default() *Forward {
	return defaultForward
}

type element struct {
	cc      *grpc.ClientConn
	gateway *runtime.ServeMux
}
type Forward struct {
	web  web.Helper
	keys map[int64]element
	rw   sync.RWMutex
}

func newForward() *Forward {
	return &Forward{
		keys: make(map[int64]element),
	}
}
func (f *Forward) Del(id int64) {
	f.rw.Lock()
	defer f.rw.Unlock()
	_, exists := f.keys[id]
	if exists {
		delete(f.keys, id)
	}
}
func (f *Forward) Put(id int64, cc *grpc.ClientConn, gateway *runtime.ServeMux) {
	f.rw.Lock()
	defer f.rw.Unlock()

	old, exists := f.keys[id]
	if exists {
		old.cc.Close()
	}
	f.keys[id] = element{
		cc:      cc,
		gateway: gateway,
	}
}
func (f *Forward) Forward(id int64, c *gin.Context) {
	f.rw.RLock()
	ele, exists := f.keys[id]
	f.rw.RUnlock()

	if !exists {
		f.web.Error(c, http.StatusNotFound, codes.NotFound, `slave id not found: `+strconv.FormatInt(id, 10))
		return
	}

	// token

	// ServeHTTP
	ele.gateway.ServeHTTP(c.Writer, c.Request)
}
