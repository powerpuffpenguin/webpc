package v1

import (
	"errors"
	"net/http"
	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/web"
	
	"google.golang.org/grpc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/powerpuffpenguin/sessionid"
)

type Logger struct {
	web.Helper
}

func (h Logger) Register(cc *grpc.ClientConn, router *gin.RouterGroup) {
	r := router.Group(`logger`)
	r.GET(`attach`, h.attach)
}
func (h Logger) checkRoot(c *gin.Context) (code int, msg string) {
	userdata, e := h.ShouldBindUserdata(c)
	if e != nil {
		if sessionid.IsTokenExpired(e) {
			code = http.StatusUnauthorized
		} else if errors.Is(e, sessionid.ErrTokenNotExists) {
			code = http.StatusForbidden
		} else {
			code = http.StatusInternalServerError
		}
		msg = e.Error()
		return
	}
	if !userdata.AuthAny(db.Root) {
		code = http.StatusForbidden
		msg = `permission denied`
		return
	}
	return
}
func (h Logger) attach(c *gin.Context) {
	ws, e := h.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	defer ws.Close()
	code, msg := h.checkRoot(c)
	if code != 0 {
		h.WSWriteClose(ws, uint16(code), msg)
		return
	}

	done := make(chan struct{})
	listener := logger.NewSnapshotListener(done)
	logger.AddListener(listener)
	go h.readWS(ws, done)
	var (
		ch      = listener.Channel()
		working = true
		data    []byte
	)
	for working {
		select {
		case <-done:
			working = false
		case data = <-ch:
			if len(data) > 0 {
				e = ws.WriteMessage(websocket.TextMessage, data)
				if e != nil {
					working = false
				}
			}
		}
	}
	logger.RemoveListener(listener)
}
func (h Logger) readWS(ws *websocket.Conn, done chan<- struct{}) {
	var e error
	for {
		_, _, e = ws.ReadMessage()
		if e != nil {
			break
		}
	}
	close(done)
}
