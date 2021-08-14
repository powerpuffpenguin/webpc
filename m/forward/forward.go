package forward

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/m/web"
	"github.com/powerpuffpenguin/webpc/sessions"
	signal_group "github.com/powerpuffpenguin/webpc/signal/group"
	signal_slave "github.com/powerpuffpenguin/webpc/signal/slave"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	web          web.Helper
	keys         map[int64]element
	subscription map[*Subscription]bool
	rw           sync.RWMutex
}

func newForward() *Forward {
	return &Forward{
		keys:         make(map[int64]element),
		subscription: make(map[*Subscription]bool),
	}
}
func (f *Forward) Del(id int64) {
	f.rw.Lock()
	defer f.rw.Unlock()
	_, exists := f.keys[id]
	if exists {
		delete(f.keys, id)

		for s := range f.subscription {
			s.put(State{
				ID:    id,
				Ready: false,
			})
		}
	}
}
func (f *Forward) Put(id int64, cc *grpc.ClientConn, gateway *runtime.ServeMux) {
	f.rw.Lock()
	defer f.rw.Unlock()
	for s := range f.subscription {
		s.put(State{
			ID:    id,
			Ready: true,
		})
	}

	old, exists := f.keys[id]
	if exists {
		old.cc.Close()
	}
	f.keys[id] = element{
		cc:      cc,
		gateway: gateway,
	}
}
func (f *Forward) Get(id int64) (cc *grpc.ClientConn, e error) {
	f.rw.RLock()
	ele, exists := f.keys[id]
	if exists {
		cc = ele.cc
	} else {
		e = status.Error(codes.NotFound, `slave id not found: `+strconv.FormatInt(id, 10))
	}
	f.rw.RUnlock()
	return
}
func (f *Forward) forwardShared(gateway *runtime.ServeMux, c *gin.Context, expired bool) {
	if expired {
		c.Request.Header.Set(`Authorization`, `Bearer Expired`)
	} else {
		c.Request.Header.Del(`Authorization`)
	}
	gateway.ServeHTTP(c.Writer, c.Request)
}
func (f *Forward) Forward(id int64, c *gin.Context) {
	f.rw.RLock()
	ele, exists := f.keys[id]
	f.rw.RUnlock()
	if !exists {
		f.web.Error(c, status.Error(codes.NotFound, `slave id not found: `+strconv.FormatInt(id, 10)))
		return
	}
	// public api
	shared := false
	if c.Request.URL.Path == `/api/forward/v1/fs/` || strings.HasPrefix(c.Request.URL.Path, `/api/forward/v1/fs/`) {
		shared = true
	}

	// userdata
	userdata, e := f.web.ShouldBindUserdata(c)
	if e != nil {
		if shared {
			f.forwardShared(ele.gateway, c, status.Code(e) == codes.Unauthenticated)
		} else {
			f.web.Error(c, e)
		}
		return
	}

	// check group
	e = f.checkGroup(c, id, &userdata)
	if e == nil {
		// token
		b, e := json.Marshal(userdata)
		if e == nil {
			c.Request.Header.Set(`Authorization`, `Bearer `+base64.RawURLEncoding.EncodeToString(b))
		} else if shared {
			f.forwardShared(ele.gateway, c, true)
			return
		} else {
			f.web.Error(c, e)
			return
		}
	} else {
		if shared {
			f.forwardShared(ele.gateway, c, status.Code(e) == codes.Unavailable)
		} else {
			f.web.Error(c, e)
		}
		return
	}

	// ServeHTTP
	ele.gateway.ServeHTTP(c.Writer, c.Request)
}

func (f *Forward) checkGroup(c *gin.Context, id int64, userdata *sessions.Userdata) (e error) {
	if id == 0 {
		if !userdata.AuthAny(db.Root, db.Server) {
			e = status.Error(codes.PermissionDenied, `permission denied`)
			return
		}
	} else {
		if !userdata.AuthAny(db.Root) && userdata.Parent != 1 {
			parents, err := signal_group.IDS(c.Request.Context(), userdata.Parent, false)
			if err != nil {
				e = err
				return
			}
			bean, err := signal_slave.Get(c.Request.Context(), id)
			if err != nil {
				e = err
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
				e = status.Error(codes.PermissionDenied, `permission denied`)
				return
			}
		}
	}
	return
}
func (f *Forward) Subscribe(ctx context.Context) (s *Subscription) {
	s = &Subscription{
		forward: f,
		ctx:     ctx,
		ch:      make(chan State, 100),
		keys:    make(map[int64]bool),
	}
	f.rw.Lock()
	f.subscription[s] = true
	f.rw.Unlock()
	return
}
func (f *Forward) change(s *Subscription, targets []int64) {
	f.rw.RLock()
	for k := range s.keys {
		delete(s.keys, k)
	}
	for _, id := range targets {
		s.keys[id] = true
	}
	for id := range s.keys {
		_, exists := f.keys[id]
		s.put(State{
			ID:    id,
			Ready: exists,
		})
	}
	f.rw.RUnlock()
}
func (f *Forward) Unsubscribe(s *Subscription) (exists bool) {
	f.rw.Lock()
	if _, exists = f.subscription[s]; exists {
		delete(f.subscription, s)
	}
	f.rw.Unlock()
	return
}
