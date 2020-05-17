package api

import "github.com/dxvgef/tsing"

// 设置路由
func SetRouter(engine *tsing.Engine) {
	var dataHandler Data
	engine.POST("/data/load-all", dataHandler.LoadAll)
	engine.POST("/data/save-all", dataHandler.SaveAll)
}
