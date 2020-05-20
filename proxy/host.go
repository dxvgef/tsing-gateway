package proxy

import (
	"strings"
)

// 写入主机，如果存在则覆盖，不存在则创建
func (p *Engine) SetHost(hostname, routeGroupID string) error {
	hostname = strings.ToLower(hostname)
	p.Hosts[hostname] = routeGroupID
	return nil
}

// 删除主机
func (p *Engine) DelHost(hostname string) error {
	hostname = strings.ToLower(hostname)
	delete(p.Hosts, hostname)
	return nil
}

// 匹配主机名，返回对应的路由组ID
func (p *Engine) matchHost(reqHost string) (string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if _, exist := p.Hosts[reqHost]; exist {
		return p.Hosts[reqHost], true
	}
	reqHost = "*"
	if _, exist := p.Hosts[reqHost]; exist {
		return p.Hosts[reqHost], true
	}
	return p.Hosts[reqHost], false
}
