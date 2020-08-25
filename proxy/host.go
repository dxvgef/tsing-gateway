package proxy

import (
	"strings"

	"github.com/rs/zerolog/log"

	"local/global"
	"local/middleware"
)

// 写入主机，如果存在则覆盖，不存在则创建
func SetHost(hostname string, host global.HostType) error {
	hostname = strings.ToLower(hostname)
	global.Hosts.Store(hostname, host)

	// 更新中间件
	mwLen := len(host.Middleware)
	if mwLen == 0 {
		global.HostMiddleware.Delete(hostname)
		return nil
	}
	mw := make([]global.MiddlewareType, mwLen)
	for k := range host.Middleware {
		m, err := middleware.Build(host.Middleware[k].Name, host.Middleware[k].Config, false)
		if err != nil {
			log.Err(err).Caller().Send()
			return err
		}
		mw = append(mw, m)
	}
	global.HostMiddleware.Store(hostname, mw)
	return nil
}

// 删除主机
func DelHost(hostname string) error {
	hostname = strings.ToLower(hostname)
	global.Hosts.Delete(hostname)
	global.HostMiddleware.Delete(hostname)
	return nil
}

// 匹配主机名
func matchHost(hostname string) global.HostType {
	pos := strings.LastIndex(hostname, ":")
	if pos > -1 {
		hostname = hostname[:pos]
	}
	if v, exist := global.Hosts.Load(hostname); exist {
		host, ok := v.(global.HostType)
		if !ok {
			log.Error().Caller().Msg("类型断言失败")
			return host
		}
		return host
	}
	hostname = "*"
	if v, exist := global.Hosts.Load(hostname); exist {
		host, ok := v.(global.HostType)
		if !ok {
			log.Error().Caller().Msg("类型断言失败")
			return host
		}
		return host
	}

	v, exist := global.Hosts.Load(hostname)
	if !exist {
		return global.HostType{}
	}
	host, ok := v.(global.HostType)
	if !ok {
		log.Error().Caller().Msg("类型断言失败")
		return host
	}
	return host
}
