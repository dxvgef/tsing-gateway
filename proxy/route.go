package proxy

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

// 设置路由组及路由，如果存在则更新，不存在则新建
func SetRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
	if routeGroupID == "" {
		return errors.New("路由组ID不能为空")
	}
	if routePath == "" {
		return errors.New("路径不能为空")
	}
	if routeMethod == "" {
		return errors.New("HTTP方法不能为空")
	}
	if !global.InStr(global.Methods, routeMethod) {
		return errors.New("HTTP方法无效")
	}
	if _, exist := global.Routes[routeGroupID]; !exist {
		global.Routes[routeGroupID] = make(map[string]map[string]string)
	}
	if _, exist := global.Routes[routeGroupID][routePath]; !exist {
		global.Routes[routeGroupID][routePath] = make(map[string]string)
	}
	global.Routes[routeGroupID][routePath][routeMethod] = upstreamID
	return nil
}

// 删除路由
func DelRoute(routeGroupID, routePath, routeMethod string) error {
	if routeGroupID == "" {
		routeGroupID = global.SnowflakeNode.Generate().String()
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
	delete(global.Routes[routeGroupID][routePath], routeMethod)
	return nil
}

// 匹配路由，返回集群ID和匹配结果的HTTP状态码
func matchRoute(req *http.Request) (upstream global.UpstreamType, status int) {
	routeGroupID := ""
	reqPath := req.URL.Path
	reqMethod := req.Method
	matchResult := false

	// 匹配主机
	routeGroupID, matchResult = matchHost(req.Host)
	if !matchResult {
		status = http.StatusServiceUnavailable
		return
	}
	// 匹配路径
	reqPath, matchResult = matchPath(routeGroupID, reqPath)
	if !matchResult {
		status = http.StatusNotFound
		return
	}
	// 匹配方法
	reqMethod, matchResult = matchMethod(routeGroupID, reqPath, reqMethod)
	if !matchResult {
		status = http.StatusMethodNotAllowed
		return
	}
	// 匹配上游
	upstreamID := global.Routes[routeGroupID][reqPath][reqMethod]
	upstream, matchResult = matchUpstream(upstreamID)
	if !matchResult {
		status = http.StatusNotImplemented
		return
	}
	return
}

// 匹配路径，返回最终匹配到的路径
func matchPath(routeGroupID, routePath string) (string, bool) {
	if routePath == "" {
		routePath = "/"
	}
	// 先尝试完全匹配路径
	if _, exist := global.Routes[routeGroupID][routePath]; exist {
		return routePath, true
	}
	// 尝试模糊匹配
	pathLastChar := routePath[len(routePath)-1]
	if pathLastChar != 47 {
		pos := strings.LastIndex(routePath, path.Base(routePath))
		routePath = routePath[:pos]
	}
	routePath = routePath + "*"
	if _, exist := global.Routes[routeGroupID][routePath]; exist {
		return routePath, true
	}
	return routePath, false
}

// 匹配方法，返回对应的集群ID
func matchMethod(routeGroupID, routePath, routeMethod string) (string, bool) {
	if _, exist := global.Routes[routeGroupID][routePath][routeMethod]; exist {
		return routeMethod, true
	}
	routeMethod = "ANY"
	if _, exist := global.Routes[routeGroupID][routePath][routeMethod]; exist {
		return routeMethod, true
	}
	return routeMethod, false
}
