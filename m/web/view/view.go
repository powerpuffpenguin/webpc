package view

import (
	"net/http"
	"strings"

	"github.com/powerpuffpenguin/webpc/m/web"
	"github.com/powerpuffpenguin/webpc/static"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
)

const BaseURL = `view`

var locales = map[string]http.FileSystem{
	`zh-Hant`: static.ZhHant(),
	`zh-Hans`: static.ZhHans(),
	`en-US`:   static.EnglishUS(),
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
			if strings.Contains(locale, `cn`) || strings.Contains(locale, `hans`) {
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
	router.GET(``, h.redirect)
	router.HEAD(``, h.redirect)
	router.GET(`index`, h.redirect)
	router.HEAD(`index`, h.redirect)
	router.GET(`index.html`, h.redirect)
	router.HEAD(`index.html`, h.redirect)
	router.GET(`view`, h.redirect)
	router.HEAD(`view`, h.redirect)
	router.GET(`view/`, h.redirect)
	router.HEAD(`view/`, h.redirect)

	r := router.Group(BaseURL)
	r.Use(h.Compression())
	r.GET(`:locale`, h.viewOrRedirect)
	r.HEAD(`:locale`, h.viewOrRedirect)
	r.GET(`:locale/*path`, h.view)
	r.HEAD(`:locale/*path`, h.view)
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
		h.Error(c, status.Error(codes.NotFound, `not support locale`))
	}
}
