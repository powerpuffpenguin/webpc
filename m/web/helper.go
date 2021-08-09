package web

import (
	"encoding/binary"
	"net/http"

	"github.com/powerpuffpenguin/webpc/m/web/contrib/compression"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 32,
	WriteBufferSize: 1024 * 32,
}

var Offered = []string{
	binding.MIMEJSON,
}
var _compression = compression.Compression(
	compression.BrDefaultCompression,
	compression.GzDefaultCompression,
)

type Helper int

func (h Helper) NegotiateData(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func (h Helper) BindURI(c *gin.Context, obj interface{}) (e error) {
	e = c.ShouldBindUri(obj)
	if e != nil {
		h.Error(c, http.StatusBadRequest, codes.InvalidArgument, e.Error())
		return
	}
	return
}

func (h Helper) Error(c *gin.Context, code int, gcode codes.Code, msg string) {
	c.JSON(code, gin.H{
		`code`:    gcode,
		`message`: msg,
	})
}

func (h Helper) Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return h.BindWith(c, obj, b)
}

func (h Helper) BindWith(c *gin.Context, obj interface{}, b binding.Binding) (e error) {
	e = c.ShouldBindWith(obj, b)
	if e != nil {
		h.Error(c, http.StatusBadRequest, codes.InvalidArgument, e.Error())
		return
	}
	return
}

func (h Helper) BindQuery(c *gin.Context, obj interface{}) error {
	return h.BindWith(c, obj, binding.Query)
}

func (h Helper) CheckWebsocket(c *gin.Context) {
	if !c.IsWebsocket() {
		c.Abort()
		return
	}
}

func (h Helper) Compression() gin.HandlerFunc {
	return _compression
}

func (h Helper) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, responseHeader)
}

func (h Helper) WSWriteClose(ws *websocket.Conn, code uint16, msg string) (e error) {
	b := make([]byte, 2+len(msg))
	binary.BigEndian.PutUint16(b, code)
	copy(b[2:], msg)
	e = ws.WriteMessage(websocket.CloseMessage, b)
	return
}
func (h Helper) WSWriteError(ws *websocket.Conn, e error) {
	ws.WriteJSON(gin.H{
		`code`: status.Code(e),
		`emsg`: e.Error(),
	})
}
func (h Helper) WSForward(ws *websocket.Conn, f Forward) {
	go h.wsRequest(ws, f)
	for {
		e := f.Response()
		if e != nil {
			h.WSWriteError(ws, e)
			break
		}
	}
	ws.Close()
}
func (h Helper) wsRequest(ws *websocket.Conn, f Forward) {
	for {
		t, p, e := ws.ReadMessage()
		if e != nil {
			h.WSWriteError(ws, e)
			break
		}
		e = f.Request(t, p)
		if e != nil {
			h.WSWriteError(ws, e)
			break
		}
	}
	ws.Close()
	f.CloseSend()
}
