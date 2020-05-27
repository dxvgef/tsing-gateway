package proxy

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

// 设置全局中间件，同时构建并更新所有全局中间件的实例
func SetGlobalMiddleware(config string) error {
	if config == "" {
		global.GlobalMiddleware = nil
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
		global.GlobalMiddleware = nil
		return nil
	}
	mw := make([]global.MiddlewareType, mwConfigLen)
	// 根据配置构建中间件实例
	for k := range mwConfig {
		log.Debug().Caller().Interface("mwConfig", mwConfig[k]).Send()
		m, err = middleware.Build(mwConfig[k].Name, mwConfig[k].Config, false)
		if err != nil {
			return err
		}
		mw = append(mw, m)
	}
	global.GlobalMiddleware = mw
	return nil
}

// 构建并更新某个上游的所有中间件的实例
func SetUpstreamMiddleware(upstreamID string, mwConfig []global.ModuleConfig) error {
	if upstreamID == "" {
		return errors.New("上游ID不能为空")
	}
	mwConfigLen := len(mwConfig)
	if mwConfigLen == 0 {
		global.UpstreamMiddleware[upstreamID] = nil
		return nil
	}
	var (
		err error
		m   global.MiddlewareType
	)

	mw := make([]global.MiddlewareType, mwConfigLen)
	// 根据配置构建中间件实例
	for k := range mwConfig {
		m, err = middleware.Build(mwConfig[k].Name, mwConfig[k].Config, false)
		if err != nil {
			return err
		}
		mw[k] = m
	}
	global.UpstreamMiddleware[upstreamID] = mw
	return nil
}
