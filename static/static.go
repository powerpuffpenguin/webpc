package static

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed favicon.ico
var favicon embed.FS

func Favicon(c *gin.Context) {
	c.Header("Cache-Control", "max-age=2419200")
	c.FileFromFS(`favicon.ico`, http.FS(favicon))
}

//go:embed document/*
var document embed.FS

func Document() http.FileSystem {
	f, e := fs.Sub(document, `document`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}

//go:embed public/*
var static embed.FS

func Static() http.FileSystem {
	f, e := fs.Sub(static, `public`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}

//go:embed view/*
var view embed.FS

func ZhHant() http.FileSystem {
	f, e := fs.Sub(view, `view/zh-Hant`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}
func ZhHans() http.FileSystem {
	f, e := fs.Sub(view, `view/zh-Hans`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}
func EnglishUS() http.FileSystem {
	f, e := fs.Sub(view, `view/en-US`)
	if e != nil {
		panic(e)
	}
	return http.FS(f)
}