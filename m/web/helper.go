package web

import (
	"net/http"

	"github.com/powerpuffpenguin/webpc/m/web/contrib/compression"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Error struct {
	Code    codes.Code `json:"code,omitempty"`
	Message string     `json:"message,omitempty"`
}

var _compression = compression.Compression(
	compression.BrDefaultCompression,
	compression.GzDefaultCompression,
)

type Helper int

func (h Helper) Response(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}
func (h Helper) Error(c *gin.Context, e error) {
	if e == nil {
		h.ResponseError(c, Error{
			Code:    codes.OK,
			Message: codes.OK.String(),
		})
		return
	}
	h.ResponseError(c, Error{
		Code:    status.Code(e),
		Message: e.Error(),
	})
}
func (h Helper) ResponseError(c *gin.Context, err Error) {
	code := http.StatusOK
	switch err.Code {
	case codes.OK:
		// code = http.StatusOK
	case codes.Canceled:
		code = http.StatusRequestTimeout
	case codes.Unknown:
		code = http.StatusInternalServerError
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.DeadlineExceeded:
		code = http.StatusGatewayTimeout
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.PermissionDenied:
		code = http.StatusForbidden
	case codes.ResourceExhausted:
		code = http.StatusTooManyRequests
	case codes.FailedPrecondition:
		code = http.StatusBadRequest
	case codes.Aborted:
		code = http.StatusConflict
	case codes.OutOfRange:
		code = http.StatusBadRequest
	case codes.Unimplemented:
		code = http.StatusNotImplemented
	case codes.Internal:
		code = http.StatusInternalServerError
	case codes.Unavailable:
		code = http.StatusServiceUnavailable
	case codes.DataLoss:
		code = http.StatusInternalServerError
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}
	c.JSON(code, err)
}
func (h Helper) BindURI(c *gin.Context, obj interface{}) (e error) {
	e = c.ShouldBindUri(obj)
	if e != nil {
		e = status.Error(codes.InvalidArgument, e.Error())
		h.Error(c, e)
		return
	}
	return
}
func (h Helper) Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return h.BindWith(c, obj, b)
}
func (h Helper) BindWith(c *gin.Context, obj interface{}, b binding.Binding) (e error) {
	e = c.ShouldBindWith(obj, b)
	if e != nil {
		e = status.Error(codes.InvalidArgument, e.Error())
		h.Error(c, e)
		return
	}
	return
}
func (h Helper) BindQuery(c *gin.Context, obj interface{}) error {
	return h.BindWith(c, obj, binding.Query)
}
func (h Helper) Compression() gin.HandlerFunc {
	return _compression
}
