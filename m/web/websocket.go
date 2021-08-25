package web

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32,
	WriteBufferSize: 1024 * 32,
}

func (h Helper) NewContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	token := h.GetToken(c)
	if token != `` {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
			`Authorization`, token,
		))
	}
	return ctx
}

func (h Helper) CheckWebsocket(c *gin.Context) {
	if !c.IsWebsocket() {
		h.Error(c, status.Error(codes.InvalidArgument, `expect websocket`))
		c.Abort()
		return
	}
}
func (h Helper) Websocket(c *gin.Context, responseHeader http.Header) (conn Websocket, e error) {
	if !c.IsWebsocket() {
		e = status.Error(codes.InvalidArgument, `expect websocket`)
		h.Error(c, e)
		return
	}
	ws, e := upgrader.Upgrade(c.Writer, c.Request, responseHeader)
	if e != nil {
		e = status.Error(codes.Unknown, e.Error())
		h.Error(c, e)
		return
	}
	conn = Websocket{ws}
	return
}

type Websocket struct {
	*websocket.Conn
}

func (w Websocket) SendMessage(m proto.Message) error {
	b, e := Marshal(m)
	if e != nil {
		return e
	}
	return w.WriteMessage(websocket.TextMessage, b)
}
func (w Websocket) SendBinary(b []byte) error {
	return w.WriteMessage(websocket.BinaryMessage, b)
}
func (w Websocket) Send(v interface{}) error {
	return w.WriteJSON(v)
}
func (w Websocket) Success() error {
	return w.Send(Error{
		Code:    codes.OK,
		Message: codes.OK.String(),
	})
}
func (w Websocket) Error(e error) error {
	if e == nil {
		return w.Send(Error{
			Code:    codes.OK,
			Message: codes.OK.String(),
		})
	} else {
		return w.Send(Error{
			Code:    status.Code(e),
			Message: e.Error(),
		})
	}
}
func (w Websocket) Forward(f Forward) {
	work := newWebsocketForward(w, f)
	work.Serve()
}

type websocketForward struct {
	w      Websocket
	f      Forward
	closed int32
	cancel chan struct{}
}

func newWebsocketForward(w Websocket, f Forward) *websocketForward {
	return &websocketForward{
		w:      w,
		f:      f,
		cancel: make(chan struct{}),
	}
}
func (wf *websocketForward) Serve() {
	go wf.request()
	go wf.response()
	<-wf.cancel
	wf.w.Close()
	wf.f.CloseSend()
}
func (wf *websocketForward) request() {
	var counted uint64
	for {
		t, p, e := wf.w.ReadMessage()
		if e != nil {
			break
		}
		e = wf.f.Request(counted, t, p)
		if e != nil {
			wf.w.Error(e)
			break
		}
		counted++
	}

	if wf.closed == 0 &&
		atomic.SwapInt32(&wf.closed, 1) == 0 {
		close(wf.cancel)
	}
}
func (wf *websocketForward) response() {
	var counted uint64
	for {
		e := wf.f.Response(counted)
		if e != nil {
			wf.w.Error(e)
			break
		}
		counted++
	}
	if wf.closed == 0 &&
		atomic.SwapInt32(&wf.closed, 1) == 0 {
		close(wf.cancel)
	}
}
