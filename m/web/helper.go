package web

import (
	"encoding/binary"
	"net/http"
	"github.com/powerpuffpenguin/webpc/m/web/contrib/compression"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEXML,
	binding.MIMEYAML,
}
var _compression = compression.Compression(
	compression.BrDefaultCompression,
	compression.GzDefaultCompression,
)

type Helper int

func (h Helper) NegotiateData(c *gin.Context, code int, data interface{}) {
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEXML:
		c.XML(code, data)
	case binding.MIMEYAML:
		c.YAML(code, data)
	default:
		// default use json
		c.JSON(code, data)
	}
}

func (h Helper) BindURI(c *gin.Context, obj interface{}) (e error) {
	e = c.ShouldBindUri(obj)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	return
}

func (h Helper) NegotiateError(c *gin.Context, code int, e error) {
	c.String(code, e.Error())
}

func (h Helper) NegotiateErrorString(c *gin.Context, code int, e string) {
	c.String(code, e)
}

func (h Helper) Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return h.BindWith(c, obj, b)
}

func (h Helper) BindWith(c *gin.Context, obj interface{}, b binding.Binding) (e error) {
	e = c.ShouldBindWith(obj, b)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
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
