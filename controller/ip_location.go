package controller

import (
	"github.com/vicanso/cod"
	"github.com/vicanso/location/router"
	"github.com/vicanso/location/service"
)

type (
	// ipLocationCtrl ip location ctrl
	ipLocationCtrl struct{}
)

var (
	unknownIPLocation = &service.IPLocation{}
)

func init() {
	g := router.NewGroup("/ip-locations")
	ctrl := ipLocationCtrl{}
	g.GET("/json/:ip", ctrl.getLocation)
}

func (ctrl ipLocationCtrl) getLocation(c *cod.Context) (err error) {
	c.NoCache()
	ipAddr := c.Param("ip")
	// 如果为此地址，则使用客户端IP
	if ipAddr == "127.0.0.1" {
		ipAddr = c.RealIP()
	}
	ip, err := service.ConvertIP2Uint32(ipAddr)
	if err != nil {
		return
	}
	location := service.GetLocationByIP(ip)
	if location == nil {
		location = unknownIPLocation
	}
	location.IP = ipAddr
	c.Body = location
	return
}
