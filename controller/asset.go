package controller

import (
	"bytes"
	"os"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/vicanso/cod"
	"github.com/vicanso/cod/middleware"
	"github.com/vicanso/location/router"
)

type (
	// assetCtrl asset ctrl
	assetCtrl struct {
	}
	staticFile struct {
		box *packr.Box
	}
)

var (
	box = packr.New("asset", "../web/build")
)

func (sf *staticFile) Exists(file string) bool {
	return sf.box.Has(file)
}
func (sf *staticFile) Get(file string) ([]byte, error) {
	return sf.box.Find(file)
}
func (sf *staticFile) Stat(file string) os.FileInfo {
	return nil
}

func init() {
	g := router.NewGroup("")
	ctrl := assetCtrl{}
	g.GET("/", noQuery, ctrl.index)
	g.GET("/favicon.ico", ctrl.favIcon)

	sf := &staticFile{
		box: box,
	}
	g.GET("/static/*file", middleware.NewStaticServe(sf, middleware.StaticServeConfig{
		// 客户端缓存一年
		MaxAge: 365 * 24 * 3600,
		// 缓存服务器缓存一个小时
		SMaxAge:             60 * 60,
		DenyQueryString:     true,
		DisableLastModified: true,
	}))
}

func sendFile(c *cod.Context, file string) (err error) {
	buf, err := box.Find(file)
	if err != nil {
		return
	}
	c.SetFileContentType(file)
	c.BodyBuffer = bytes.NewBuffer(buf)
	return
}

func (ctrl assetCtrl) index(c *cod.Context) (err error) {
	c.CacheMaxAge("10s")
	return sendFile(c, "index.html")
}

func (ctrl assetCtrl) favIcon(c *cod.Context) (err error) {
	c.CacheMaxAge("10m")
	return sendFile(c, "images/favicon.ico")
}
