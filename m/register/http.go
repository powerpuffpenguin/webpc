package register

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/forward"
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
			if query.AccessToken == `` {
				query.AccessToken, _ = c.Cookie(helper.CookieName)
				if query.AccessToken != `` {
					c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
				}
			} else {
				c.Request.Header.Set(`Authorization`, `Bearer `+query.AccessToken)
			}
		}
		if strings.HasPrefix(c.Request.URL.Path, `/api/forward/`) {
			var query struct {
				ID string `form:"slave_id" binding:"required"`
			}
			e := c.BindQuery(&query)
			if e != nil {
				return
			}
			forward.Default().Forward(query.ID, c)
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
	makeStatic(engine)
	var views view.Helper
	views.Register(&engine.RouterGroup)
	// other gin route
	var apis api.Helper
	apis.Register(cc, &engine.RouterGroup)
}

type Static struct {
	web.Helper
	fileSystem http.FileSystem
}

func makeStatic(engine *gin.Engine) {
	fileSystem, e := fs.NewWithNamespace(`static`)
	if e != nil {
		logger.Logger.Panic(`statik static error`,
			zap.Error(e),
		)
	}
	h := Static{
		fileSystem: fileSystem,
	}
	compression := h.Compression()
	engine.GET(`favicon.ico`, compression, h.favicon)
	engine.HEAD(`favicon.ico`, compression, h.favicon)
	engine.Group(`static`).Use(compression).StaticFS(``, fileSystem)
}
func (h Static) favicon(c *gin.Context) {
	h.NegotiateFilesystem(c, h.fileSystem, `/favicon.ico`, false)
}
