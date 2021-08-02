package view

import (
	"net/http"
	"os"
	"strings"

	"github.com/powerpuffpenguin/webpc/logger"
	"github.com/powerpuffpenguin/webpc/m/web"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"go.uber.org/zap"
)

const BaseURL = `view`

var locales = make(map[string]http.FileSystem)
var supported = []string{
	`zh-Hant`,
	`zh-Hans`,
	`en-US`,
}

func localeResolution(accept string) string {
	var (
		index  int
		locale string
	)
	for accept != `` {
		index = strings.IndexAny(accept, `,;`)
		if index == -1 {
			locale = strings.ToLower(strings.TrimSpace(accept))
			accept = ``
		} else {
			locale = strings.ToLower(strings.TrimSpace(accept[:index]))
			accept = accept[index+1:]
		}

		if locale == `zh` {
			return `zh-Hant`
		} else if locale == `en` {
			return `en-US`
		} else if strings.HasPrefix(locale, `zh-`) {
			if strings.Index(locale, `cn`) != -1 || strings.Index(locale, `hans`) != -1 {
				return `zh-Hans`
			}
			return `zh-Hant`
		} else if strings.HasPrefix(locale, `en-`) {
			return `en-US`
		}
	}
	return `en-US`
}

type Helper struct {
	web.Helper
}

func (h Helper) Register(router *gin.RouterGroup) {
	for _, key := range supported {
		locale, e := fs.NewWithNamespace(key)
		if e != nil {
			if ce := logger.Logger.Check(zap.FatalLevel, `New FileSystem error`); ce != nil {
				ce.Write(
					zap.Error(e),
					zap.String(`namespace`, key),
				)
			}
			os.Exit(1)
		}
		locales[key] = locale
	}

	router.GET(``, h.redirect)
	router.GET(`index`, h.redirect)
	router.GET(`index.html`, h.redirect)
	router.GET(`view`, h.redirect)
	router.GET(`view/`, h.redirect)

	r := router.Group(BaseURL)
	r.Use(h.Compression())
	r.GET(`:locale`, h.viewOrRedirect)
	r.GET(`:locale/*path`, h.view)
}

func (h Helper) redirectAngular(c *gin.Context) {
	var obj struct {
		Path string `uri:"path"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	request := c.Request
	str := strings.ToLower(strings.TrimSpace(request.Header.Get(`Accept-Language`)))
	strs := strings.Split(str, `;`)
	str = strings.TrimSpace(strs[0])
	strs = strings.Split(str, `,`)
	str = strings.TrimSpace(strs[0])
	if strings.HasPrefix(str, `zh-`) {
		if strings.Index(str, `cn`) != -1 || strings.Index(str, `hans`) != -1 {
			c.Redirect(http.StatusFound, `/view/zh-Hans/`+obj.Path)
			return
		}
		c.Redirect(http.StatusFound, `/view/zh-Hant/`+obj.Path)
		return
	}
	c.Redirect(http.StatusFound, `/view/en-US/`+obj.Path)
}
func (h Helper) redirect(c *gin.Context) {
	c.Redirect(http.StatusFound, `/view/`+localeResolution(c.Request.Header.Get(`Accept-Language`))+`/`)
}
func (h Helper) viewOrRedirect(c *gin.Context) {
	var obj struct {
		Locale string `uri:"locale"`
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	if _, ok := locales[obj.Locale]; ok {
		c.Redirect(http.StatusFound, `/view/`+obj.Locale+`/`)
	} else {
		h.redirect(c)
	}
}
func (h Helper) view(c *gin.Context) {
	var obj struct {
		Locale string `uri:"locale" binding:"required"`
		Path   string `uri:"path" `
	}
	e := h.BindURI(c, &obj)
	if e != nil {
		return
	}
	if f, ok := locales[obj.Locale]; ok {
		c.Header("Cache-Control", "max-age=2419200")
		h.NegotiateFilesystem(c, f, obj.Path, true)
	} else {
		h.NegotiateErrorString(c, http.StatusNotFound, `not support locale`)
	}
}
