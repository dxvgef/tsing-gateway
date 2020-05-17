package proxy

import (
	"errors"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

type RouteGroup struct {
	ID string
	p  *Engine
}

// 新建路由组
func (p *Engine) NewRouteGroup(routeGroupID string) (routeGroup RouteGroup, err error) {
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
	routeGroup.p = p
	return
}

// 设置路由组，如果存在则更新，不存在则新建
func (p *Engine) SetRouteGroup(routeGroupID string) (routeGroup RouteGroup, err error) {
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
	routeGroup.p = p
	return
}

// 删除路由组，所有子路由都将被清空
func (p *Engine) DelRouteGroup(routeGroupID string) error {
	if routeGroupID == "" {
		return errors.New("routeGroupID不能为空")
	}
	delete(p.Routes, routeGroupID)
	return nil
}

// 在路由组内设置路由，如果存在则更新，不存在则新建
func (g *RouteGroup) SetRoute(path, method, upstreamID string) error {
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
	// if _, exist := g.p.Upstreams[upstreamID]; !exist {
	// 	return errors.New("上游ID:" + upstreamID + "不存在")
	// }
	if _, exist := g.p.Routes[g.ID]; !exist {
		g.p.Routes[g.ID] = make(map[string]map[string]string)
	}
	if _, exist := g.p.Routes[g.ID][path]; !exist {
		g.p.Routes[g.ID][path] = make(map[string]string)
	}
	g.p.Routes[g.ID][path][method] = upstreamID

	return nil
}

// 在路由组内删除路由
func (g *RouteGroup) DelRoute(path, method string) error {
	if path == "" {
		return errors.New("path不能为空")
	}
	if method == "" {
		return errors.New("method不能为空")
	}
	method = strings.ToUpper(method)
	delete(g.p.Routes[g.ID][path], method)
	return nil
}
