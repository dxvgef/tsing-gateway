package main

import (
	"errors"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

type RouteGroup struct {
	ID    string
	proxy *Proxy
}

// 新建路由组
func (p *Proxy) NewRouteGroup(routeGroupID string, persistent bool) (routeGroup RouteGroup, err error) {
	if routeGroupID == "" {
		routeGroupID = global.GetIDStr()
	}
	if routeGroupID == "" {
		err = errors.New("没有传入路由组ID,并且无法自动创建ID")
		return
	}

	if _, exist := p.Routes[routeGroupID]; exist {
		err = errors.New("路由ID:" + routeGroupID + "已存在")
		return
	}
	p.Routes[routeGroupID] = make(map[string]map[string]string)
	routeGroup.ID = routeGroupID
	routeGroup.proxy = p
	return
}

// 设置路由组，如果存在则更新，不存在则新建
func (p *Proxy) SetRouteGroup(routeGroupID string, persistent bool) (routeGroup RouteGroup, err error) {
	if routeGroupID == "" {
		routeGroupID = global.GetIDStr()
	}
	if routeGroupID == "" {
		err = errors.New("没有传入路由组ID,并且无法自动创建ID")
		return
	}
	if _, exist := p.Routes[routeGroupID]; !exist {
		p.Routes[routeGroupID] = make(map[string]map[string]string)
	}
	routeGroup.ID = routeGroupID
	routeGroup.proxy = p
	return
}

// 在路由组内设置路由，如果存在则更新，不存在则新建
func (g *RouteGroup) SetRoute(path, method, upstreamID string, persistent bool) error {
	if path == "" {
		path = "/"
	}
	if method == "" {
		method = "*"
	} else {
		method = strings.ToUpper(method)
	}
	if g.ID == "" {
		g.ID = global.GetIDStr()
	}
	if g.ID == "" {
		return errors.New("没有设置路由组ID,并且无法自动创建ID")
	}
	if _, exist := g.proxy.Upstreams[upstreamID]; !exist {
		return errors.New("上游ID:" + upstreamID + "不存在")
	}
	if _, exist := g.proxy.Routes[g.ID]; !exist {
		g.proxy.Routes[g.ID] = make(map[string]map[string]string)
	}
	if _, exist := g.proxy.Routes[g.ID][path]; !exist {
		g.proxy.Routes[g.ID][path] = make(map[string]string)
	}
	g.proxy.Routes[g.ID][path][method] = upstreamID

	return nil
}
