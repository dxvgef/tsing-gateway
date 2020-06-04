package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret
	router := engine.Group("", checkSecretFromHeader)

	var proxyHandler Proxy
	router.GET("/proxy/", proxyHandler.OutputJSON)
	router.POST("/proxy/", proxyHandler.LoadAll)
	router.PUT("/proxy/", proxyHandler.SaveAll)

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
	router.PUT("/route/:groupID/:path/:method", routeHandler.Put)
	router.DELETE("/route/:groupID/:path/:method", routeHandler.DeleteMethod)
	router.DELETE("/route/:groupID/:path/", routeHandler.DeletePath)
	router.DELETE("/route/:groupID/", routeHandler.DeleteGroup)
	router.DELETE("/route/", routeHandler.DeleteAll)
}
