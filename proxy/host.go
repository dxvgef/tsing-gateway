package proxy

import (
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
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
	global.Hosts.Delete(hostname)
	return nil
}

// 匹配主机名，返回对应的路由组ID
func matchHost(reqHost string) (string, string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if v, exist := global.Hosts.Load(reqHost); exist {
		return reqHost, v.(string), true
	}
	reqHost = "*"
	if v, exist := global.Hosts.Load(reqHost); exist {
		return reqHost, v.(string), true
	}

	v, exist := global.Hosts.Load(reqHost)
	return reqHost, v.(string), exist
}
