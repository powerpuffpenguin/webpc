package register

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/helper"
	"github.com/powerpuffpenguin/webpc/m/web"
	"github.com/powerpuffpenguin/webpc/m/web/api"
	"github.com/powerpuffpenguin/webpc/m/web/view"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func HTTP(cc *grpc.ClientConn, engine *gin.Engine, gateway *runtime.ServeMux, swagger bool) {
	var w web.Helper
	engine.NoRoute(func(c *gin.Context) {
		c.Status(http.StatusOK)
		if c.Request.Method == `GET` || c.Request.Method == `HEAD` {
			c.Request.Header.Set(`Method`, c.Request.Method)
		}
		if c.Request.Header.Get(`Authorization`) == `` {
			var query struct {
				AccessToken string `form:"access_token"`
			}
			c.ShouldBindQuery(&query)
			if query.AccessToken != `` {
				c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
			}
			if query.AccessToken == `` {
				query.AccessToken, _ = c.Cookie(helper.CookieName)
				if query.AccessToken != `` {
					c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
				}
			}
		}
		if strings.HasPrefix(c.Request.URL.Path, `/api/forward/`) {
			var query struct {
				ID int64 `form:"slave_id"`
			}
			e := c.BindQuery(&query)
			if e != nil {
				return
			}
			DefaultForward().Forward(query.ID, c)
		} else {
			gateway.ServeHTTP(c.Writer, c.Request)
		}
	})
	if swagger {
		document, e := fs.NewWithNamespace(`document`)
		if e != nil {
			logger.Logger.Panic(`statik document error`,
				zap.Error(e),
			)
		}
		r := engine.Group(`document`)
		r.Use(w.Compression())
		r.StaticFS(``, document)
	}
	static, e := fs.NewWithNamespace(`static`)
	if e != nil {
		logger.Logger.Panic(`statik static error`,
			zap.Error(e),
		)
	}
	engine.Group(`static`).Use(w.Compression()).StaticFS(``, static)
	var views view.Helper
	views.Register(&engine.RouterGroup)
	// other gin route
	var apis api.Helper
	apis.Register(cc, &engine.RouterGroup)
}
