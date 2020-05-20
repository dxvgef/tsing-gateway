package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secert
	router := engine.Group("", checkHeader)

	var dataHandler Data
	router.GET("/data/", dataHandler.LoadAll)
	router.PUT("/data/", dataHandler.SaveAll)

	var hostHandler Host
	router.POST("/host/", hostHandler.Add)
	router.PUT("/host/:hostname", hostHandler.Put)
	router.DELETE("/host/:hostname", hostHandler.Delete)

	var upstreamHandler Upstream
	router.POST("/upstream/", upstreamHandler.Add)
	router.PUT("/upstream/:id", upstreamHandler.Put)
	router.DELETE("/upstream/:id", upstreamHandler.Delete)

	var routeHandler Route
	router.POST("/route/", routeHandler.Add)
	router.PUT("/route/:key", routeHandler.Put)
	router.DELETE("/route/:key", routeHandler.Delete)
}
