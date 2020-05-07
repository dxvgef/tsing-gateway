package main

import (
	"errors"
	"net/http"
	"path"
	"strings"
)

// 新建路由组及路由
func (p *Proxy) newRoute(routeGroupID, reqPath, reqMethod, upstreamID string, persistent bool) error {
	if routeGroupID == "" {
		routeGroupID = getID()
	}
	if routeGroupID == "" {
		return errors.New("没有传入路由组ID,并且无法自动创建ID")
	}
	if reqPath == "" {
		reqPath = "/"
	}
	if reqMethod == "" {
		reqMethod = "*"
	} else {
		reqMethod = strings.ToUpper(reqMethod)
	}
	if _, exist := p.routeGroups[routeGroupID]; exist {
		return errors.New("路由组ID:" + routeGroupID + "已存在")
	}
	if _, exist := p.routeGroups[routeGroupID][reqPath]; exist {
		return errors.New("路由组ID:" + routeGroupID + "的路径:" + reqPath + "已存在")
	}
	if _, exist := p.routeGroups[routeGroupID][reqPath][reqMethod]; exist {
		return errors.New("路由组ID:" + routeGroupID + "/路径:" + reqPath + "/方法:" + reqMethod + "已存在")
	}
	if _, exist := p.upstreams[upstreamID]; !exist {
		return errors.New("上游ID:" + upstreamID + "不存在")
	}
	p.routeGroups[routeGroupID] = make(map[string]map[string]string)
	p.routeGroups[routeGroupID][reqPath] = make(map[string]string)
	p.routeGroups[routeGroupID][reqPath][reqMethod] = upstreamID
	return nil
}

// 设置路由组及路由，如果存在则更新，不存在则新建
func (p *Proxy) setRoute(routeGroupID, reqPath, reqMethod, upstreamID string, persistent bool) error {
	if routeGroupID == "" {
		routeGroupID = getID()
	}
	if routeGroupID == "" {
		return errors.New("没有传入路由组ID,并且无法自动创建ID")
	}
	if _, exist := p.upstreams[upstreamID]; !exist {
		return errors.New("上游ID:" + upstreamID + "不存在")
	}
	if reqPath == "" {
		reqPath = "/"
	}
	if reqMethod == "" {
		reqMethod = "*"
	} else {
		reqMethod = strings.ToUpper(reqMethod)
	}
	if _, exist := p.routeGroups[routeGroupID]; !exist {
		p.routeGroups[routeGroupID] = make(map[string]map[string]string)
	}
	if _, exist := p.routeGroups[routeGroupID][reqPath]; !exist {
		p.routeGroups[routeGroupID][reqPath] = make(map[string]string)
	}
	p.routeGroups[routeGroupID][reqPath][reqMethod] = upstreamID
	return nil
}

// 匹配路由，返回集群ID和匹配结果的HTTP状态码
func (p *Proxy) matchRoute(req *http.Request) (upstream Upstream, status int) {
	routeGroupID := ""
	reqPath := req.URL.Path
	reqMethod := req.Method
	matchResult := false

	// 匹配主机
	routeGroupID, matchResult = p.matchHost(req.Host)
	if !matchResult {
		status = http.StatusServiceUnavailable
		return
	}
	// 匹配路径
	reqPath, matchResult = p.matchPath(routeGroupID, reqPath)
	if !matchResult {
		status = http.StatusNotFound
		return
	}
	// 匹配方法
	reqMethod, matchResult = p.matchMethod(routeGroupID, reqPath, reqMethod)
	if !matchResult {
		status = http.StatusMethodNotAllowed
		return
	}
	// 匹配上游
	upstreamID := p.routeGroups[routeGroupID][reqPath][reqMethod]
	upstream, matchResult = p.matchUpstream(upstreamID)
	if !matchResult {
		status = http.StatusNotImplemented
		return
	}
	status = http.StatusOK
	return
}

// 匹配路径，返回最终匹配到的路径
func (p *Proxy) matchPath(routeGroupID, reqPath string) (string, bool) {
	if reqPath == "" {
		reqPath = "/"
	}
	// 先尝试完全匹配路径
	if _, exist := p.routeGroups[routeGroupID][reqPath]; exist {
		return reqPath, true
	}
	// 尝试模糊匹配
	pathLastChar := reqPath[len(reqPath)-1]
	if pathLastChar != 47 {
		pos := strings.LastIndex(reqPath, path.Base(reqPath))
		reqPath = reqPath[:pos]
	}
	reqPath = reqPath + "*"
	if _, exist := p.routeGroups[routeGroupID][reqPath]; exist {
		return reqPath, true
	}
	return reqPath, false
}

// 匹配方法，返回对应的集群ID
func (p *Proxy) matchMethod(routeGroupID, reqPath, reqMethod string) (string, bool) {
	if _, exist := p.routeGroups[routeGroupID][reqPath][reqMethod]; exist {
		return reqMethod, true
	}
	reqMethod = "*"
	if _, exist := p.routeGroups[routeGroupID][reqPath][reqMethod]; exist {
		return reqMethod, true
	}
	return reqMethod, false
}
