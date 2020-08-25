package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	// 检查secret
	router := engine.Group("", checkSecretFromHeader)

	var dataHandler Data
	router.GET("/data/", dataHandler.OutputJSON)
	router.POST("/data/", dataHandler.LoadAll)
	router.PUT("/data/", dataHandler.SaveAll)

	var hostHandler Host
	router.POST("/hosts/", hostHandler.Add)
	router.PUT("/hosts/:hostname", hostHandler.Put)
	router.DELETE("/hosts/:hostname", hostHandler.Delete)

	var serviceHandler Service
	router.POST("/services/", serviceHandler.Add)
	router.PUT("/services/:id", serviceHandler.Put)
	router.DELETE("/services/:id", serviceHandler.Delete)

	var routeHandler Route
	router.POST("/routes/", routeHandler.Add)
	router.PUT("/routes/:hostname/:path/:method", routeHandler.Put)
	router.DELETE("/routes/:hostname/:path/:method", routeHandler.DeleteMethod)
	router.DELETE("/routes/:hostname/:path/", routeHandler.DeletePath)
	router.DELETE("/routes/:hostname/", routeHandler.DeleteGroup)
	router.DELETE("/routes/", routeHandler.DeleteAll)
}
