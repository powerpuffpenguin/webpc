package web

import (
	"encoding/base64"
	"strings"

	"github.com/powerpuffpenguin/webpc/db"
	"github.com/powerpuffpenguin/webpc/m/helper"
	"github.com/powerpuffpenguin/webpc/sessionid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

type sessionValue struct {
	session *sessionid.Session
	e       error
}

func (h Helper) GetToken(c *gin.Context) string {
	return h.getToken(c)
}
func (h Helper) getToken(c *gin.Context) string {
	token := c.GetHeader(`Authorization`)
	if token != `` {
		return token
	}
	var obj struct {
		AccessToken string `form:"access_token" uri:"access_token"`
	}
	c.ShouldBindUri(&obj)
	if obj.AccessToken != `` {
		return `Bearer ` + obj.AccessToken
	}
	c.ShouldBindQuery(&obj)
	if obj.AccessToken != `` {
		return `Bearer ` + obj.AccessToken
	}
	token, _ = c.Cookie(helper.CookieName)
	if token != `` {
		return `Bearer ` + token
	}
	return ``
}

func (h Helper) accessSession(c *gin.Context) (session *sessionid.Session, e error) {
	token := h.getToken(c)
	if strings.HasPrefix(token, `Bearer `) {
		access := token[7:]
		if c.Request.URL.Path == `/api/forward/v1/fs/download_access` {
			var playdata string
			playdata, e = sessionid.DefaultManager().Verify(access)
			if e != nil {
				return
			}
			var b []byte
			b, e = base64.RawURLEncoding.DecodeString(playdata)
			if e != nil {
				return
			}
			var result sessionid.Session
			e = result.Unmarshal(b)
			if e != nil {
				return
			}
			session = &result
		} else {
			session, e = sessionid.DefaultManager().Get(c.Request.Context(), access)
		}
	}
	return
}
func (h Helper) ShouldBindSession(c *gin.Context) (session *sessionid.Session, e error) {
	v, exists := c.Get(`session`)
	if exists {
		if cache, ok := (v).(sessionValue); ok {
			session = cache.session
			e = cache.e
			return
		}
	}
	session, e = h.accessSession(c)
	if e == nil && session == nil {
		e = status.Error(codes.PermissionDenied, `token not exists`)
	}
	c.Set(`session`, sessionValue{
		session: session,
		e:       e,
	})
	return
}
func (h Helper) BindSession(c *gin.Context) (session *sessionid.Session, e error) {
	session, e = h.ShouldBindSession(c)
	if e != nil {
		h.Error(c, e)
		return
	}
	return
}
func (h Helper) ShouldBindUserdata(c *gin.Context) (userdata *sessionid.Session, e error) {
	return h.ShouldBindSession(c)
}
func (h Helper) BindUserdata(c *gin.Context) (userdata *sessionid.Session, e error) {
	return h.BindSession(c)
}

func (h Helper) CheckRoot(c *gin.Context) {
	userdata, e := h.BindUserdata(c)
	if e != nil {
		c.Abort()
		return
	}
	if userdata.AuthAny(db.Root) {
		return
	}
	h.Error(c, status.Error(codes.PermissionDenied, `permission denied`))
	c.Abort()
}
