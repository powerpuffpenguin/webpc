package forward

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/m/web"
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
	signal_slave "github.com/powerpuffpenguin/webpc/signal/slave"
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

	// userdata
	userdata, e := f.web.BindUserdata(c)
	if e != nil {
		return
	}
	b, e := json.Marshal(userdata)
	if e != nil {
		f.web.Error(c, http.StatusInternalServerError, codes.Unknown, e.Error())
		return
	}
	// check group
	if id == 0 {
		if !userdata.AuthAny(db.Root, db.Server) {
			f.web.Error(c, http.StatusForbidden, codes.PermissionDenied, `permission denied`)
			return
		}
	} else {
		if !userdata.AuthAny(db.Root) && userdata.Parent != 1 {
			parents, e := signal_group.IDS(c.Request.Context(), userdata.Parent, false)
			if e != nil {
				f.web.Error(c, http.StatusInternalServerError, codes.Unknown, e.Error())
				return
			}
			bean, e := signal_slave.Get(c.Request.Context(), id)
			if e != nil {
				f.web.Error(c, http.StatusInternalServerError, codes.Unknown, e.Error())
				return
			}
			found := false
			for _, parent := range parents.ID {
				if parent == bean.Parent {
					found = true
					break
				}
			}
			if !found {
				f.web.Error(c, http.StatusForbidden, codes.PermissionDenied, `permission denied`)
				return
			}
		}
	}
	// token
	c.Request.Header.Set(`Authorization`, base64.RawURLEncoding.EncodeToString(b))

	// ServeHTTP
	ele.gateway.ServeHTTP(c.Writer, c.Request)
}
