package proxy

import (
	"encoding/json"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

// 设置主机中间件，同时构建并更新所有主机中间件的实例
func SetHostMiddleware(hostname, config string) error {
	if hostname == "" {
		return errors.New("主机名不能为空")
	}
	if _, exist := global.Hosts.Load(hostname); !exist {
		return errors.New("主机名不存在")
	}
	if config == "" {
		global.HostMiddleware.Delete(hostname)
		return nil
	}
	var (
		err      error
		resp     = map[string]string{}
		mwConfig []global.ModuleConfig
		m        global.MiddlewareType
	)
	// 将字符串解码成模块配置
	if err = json.Unmarshal(global.StrToBytes(config), &mwConfig); err != nil {
		resp["error"] = err.Error()
		return err
	}
	mwConfigLen := len(mwConfig)
	// 如果中间件数量为0，则删除该主机的所有中间件
	if mwConfigLen == 0 {
		global.HostMiddleware.Delete(hostname)
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
	global.HostMiddleware.Store(hostname, mw)
	return nil
}

// 构建并更新某个上游的所有中间件的实例
func SetUpstreamMiddleware(upstreamID string, mwConfig []global.ModuleConfig) error {
	if upstreamID == "" {
		return errors.New("上游ID不能为空")
	}
	mwConfigLen := len(mwConfig)
	if mwConfigLen == 0 {
		global.UpstreamMiddleware.Delete(upstreamID)
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
	global.UpstreamMiddleware.Store(upstreamID, mw)
	return nil
}
