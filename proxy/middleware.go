package proxy

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

func SetMiddleware(config string) error {
	if config == "" {
		global.Middleware = nil
		return nil
	}
	var (
		err      error
		resp     = make(map[string]string)
		mwConfig []global.ModuleConfig
		m        global.MiddlewareType
	)
	// 将字符串解码成模块配置
	if err = json.Unmarshal(global.StrToBytes(config), &mwConfig); err != nil {
		resp["error"] = err.Error()
		return err
	}
	mwConfigLen := len(mwConfig)
	if mwConfigLen == 0 {
		global.Middleware = nil
		return nil
	}
	mw := make([]global.MiddlewareType, mwConfigLen)
	// 根据配置构建中间件实例
	for k := range mwConfig {
		m, err = middleware.Build(mwConfig[k].Name, mwConfig[k].Config, false)
		if err != nil {
			return err
		}
		mw = append(mw, m)
	}
	global.Middleware = mw
	return nil
}
