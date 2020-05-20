package proxy

import (
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

// 写入主机，如果存在则覆盖，不存在则创建
func SetHost(hostname, routeGroupID string) error {
	hostname = strings.ToLower(hostname)
	global.Hosts[hostname] = routeGroupID
	return nil
}

// 删除主机
func DelHost(hostname string) error {
	hostname = strings.ToLower(hostname)
	delete(global.Hosts, hostname)
	return nil
}

// 匹配主机名，返回对应的路由组ID
func matchHost(reqHost string) (string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if _, exist := global.Hosts[reqHost]; exist {
		return global.Hosts[reqHost], true
	}
	reqHost = "*"
	if _, exist := global.Hosts[reqHost]; exist {
		return global.Hosts[reqHost], true
	}
	return global.Hosts[reqHost], false
}
