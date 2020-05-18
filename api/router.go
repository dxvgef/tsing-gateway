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
	router.PUT("/host/", hostHandler.Put)
	router.DELETE("/host/:hostname", hostHandler.Del)

	var upstreamHandler Host
	router.PUT("/upstream/", upstreamHandler.Put)
	router.DELETE("/upstream/:id", upstreamHandler.Del)
}
