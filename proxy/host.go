package proxy

import (
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

// 写入主机，如果存在则覆盖，不存在则创建
func SetHost(hostname, routeGroupID string) error {
	hostname = strings.ToLower(hostname)
	global.Hosts.Store(hostname, routeGroupID)
	return nil
}

// 删除主机
func DelHost(hostname string) error {
	hostname = strings.ToLower(hostname)
	global.Hosts.Delete(hostname)
	return nil
}

// 匹配主机名，返回对应的路由组ID
func matchHost(reqHost string) (string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if v, exist := global.Hosts.Load(reqHost); exist {
		return v.(string), true
	}
	reqHost = "*"
	if v, exist := global.Hosts.Load(reqHost); exist {
		return v.(string), true
	}

	v, exist := global.Hosts.Load(reqHost)
	return v.(string), exist
}
