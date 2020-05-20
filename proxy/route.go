package proxy

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/dxvgef/gommon/slice"

	"github.com/dxvgef/tsing-gateway/global"
)

// 新建路由组及路由
func (p *Engine) NewRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
	if _, exist := p.Routes[routeGroupID][routePath][routeMethod]; exist {
		return errors.New("路由组ID:" + routeGroupID + "/路径:" + routePath + "/方法:" + routeMethod + "已存在")
	}
	return p.SetRoute(routeGroupID, routePath, routeMethod, upstreamID)
}

// 设置路由组及路由，如果存在则更新，不存在则新建
func (p *Engine) SetRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
	if routeGroupID == "" {
		return errors.New("没有传入路由组ID,并且无法自动创建ID")
	}
	// if _, exist := p.Upstreams[upstreamID]; !exist {
	// 	return errors.New("上游ID:" + upstreamID + "不存在")
	// }
	if routePath == "" {
		routePath = "/"
	}
	if routeMethod == "" {
		routeMethod = "*"
	} else {
		routeMethod = strings.ToUpper(routeMethod)
	}
	if !slice.InStr(global.Methods, routeMethod) {
		return errors.New("HTTP方法无效")
	}
	if _, exist := p.Routes[routeGroupID]; !exist {
		p.Routes[routeGroupID] = make(map[string]map[string]string)
	}
	if _, exist := p.Routes[routeGroupID][routePath]; !exist {
		p.Routes[routeGroupID][routePath] = make(map[string]string)
	}
	p.Routes[routeGroupID][routePath][routeMethod] = upstreamID
	return nil
}

// 删除路由
func (p *Engine) DelRoute(routeGroupID, routePath, routeMethod string) error {
	if routeGroupID == "" {
		routeGroupID = global.GetIDStr()
	}
	if routeGroupID == "" {
		return errors.New("routeGroupID不能为空")
	}
	if routePath == "" {
		return errors.New("reqPath不能为空")
	}
	if routeMethod == "" {
		return errors.New("reqMethod不能为空")
	}
	routeMethod = strings.ToUpper(routeMethod)
	delete(p.Routes[routeGroupID][routePath], routeMethod)
	return nil
}

// 匹配路由，返回集群ID和匹配结果的HTTP状态码
func (p *Engine) MatchRoute(req *http.Request) (upstream Upstream, status int) {
	routeGroupID := ""
	reqPath := req.URL.Path
	reqMethod := req.Method
	matchResult := false

	// 匹配主机
	routeGroupID, matchResult = p.MatchHost(req.Host)
	if !matchResult {
		status = http.StatusServiceUnavailable
		return
	}
	// 匹配路径
	reqPath, matchResult = p.MatchPath(routeGroupID, reqPath)
	if !matchResult {
		status = http.StatusNotFound
		return
	}
	// 匹配方法
	reqMethod, matchResult = p.MatchMethod(routeGroupID, reqPath, reqMethod)
	if !matchResult {
		status = http.StatusMethodNotAllowed
		return
	}
	// 匹配上游
	upstreamID := p.Routes[routeGroupID][reqPath][reqMethod]
	upstream, matchResult = p.MatchUpstream(upstreamID)
	if !matchResult {
		status = http.StatusNotImplemented
		return
	}
	status = http.StatusOK
	return
}

// 匹配路径，返回最终匹配到的路径
func (p *Engine) MatchPath(routeGroupID, reoutePath string) (string, bool) {
	if reoutePath == "" {
		reoutePath = "/"
	}
	// 先尝试完全匹配路径
	if _, exist := p.Routes[routeGroupID][reoutePath]; exist {
		return reoutePath, true
	}
	// 尝试模糊匹配
	pathLastChar := reoutePath[len(reoutePath)-1]
	if pathLastChar != 47 {
		pos := strings.LastIndex(reoutePath, path.Base(reoutePath))
		reoutePath = reoutePath[:pos]
	}
	// todo 可能要将*号换成别的符号，因为和api(tsing)的路由规则冲突
	reoutePath = reoutePath + "*"
	if _, exist := p.Routes[routeGroupID][reoutePath]; exist {
		return reoutePath, true
	}
	return reoutePath, false
}

// 匹配方法，返回对应的集群ID
func (p *Engine) MatchMethod(routeGroupID, routePath, routeMethod string) (string, bool) {
	if _, exist := p.Routes[routeGroupID][routePath][routeMethod]; exist {
		return routeMethod, true
	}
	routeMethod = "*"
	if _, exist := p.Routes[routeGroupID][routePath][routeMethod]; exist {
		return routeMethod, true
	}
	return routeMethod, false
}
