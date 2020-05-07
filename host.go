package main

import (
	"errors"
	"strings"
)

// 新建主机
func (p *Proxy) newHost(hostname, routeGroupID string, persistent bool) error {
	hostname = strings.ToLower(hostname)
	if _, ok := p.hosts[hostname]; ok {
		return errors.New("主机名:" + hostname + "已存在")
	}
	if _, exist := p.routeGroups[routeGroupID]; !exist {
		return errors.New("路由组ID:" + routeGroupID + "不存在")
	}
	p.hosts[hostname] = routeGroupID
	return nil
}

// 写入主机，如果存在则覆盖，不存在则创建
func (p *Proxy) setHost(hostname, routeGroupID string, persistent bool) error {
	hostname = strings.ToLower(hostname)
	if _, exist := p.routeGroups[routeGroupID]; !exist {
		return errors.New("路由组ID:" + routeGroupID + "不存在")
	}
	p.hosts[hostname] = routeGroupID
	return nil
}

// 匹配主机名，返回对应的路由组ID
func (p *Proxy) matchHost(reqHost string) (string, bool) {
	pos := strings.LastIndex(reqHost, ":")
	if pos > -1 {
		reqHost = reqHost[:pos]
	}
	if _, exist := p.hosts[reqHost]; exist {
		return p.hosts[reqHost], true
	}
	reqHost = "*"
	if _, exist := p.hosts[reqHost]; exist {
		return p.hosts[reqHost], true
	}
	return p.hosts[reqHost], false
}
