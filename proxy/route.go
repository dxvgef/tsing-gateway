package proxy

import (
	"errors"
	"net/http"
	"path"
	"strings"

	"local/global"
)

// 设置路由组及路由
func SetRoute(routeGroupID, routePath, routeMethod, serviceID string) error {
	if routeGroupID == "" {
		return errors.New("路由组ID不能为空")
	}
	if routePath == "" {
		return errors.New("路径不能为空")
	}
	if routeMethod == "" {
		return errors.New("HTTP方法不能为空")
	}
	if !global.InStr(global.HTTPMethods, routeMethod) {
		return errors.New("HTTP方法无效")
	}
	var key strings.Builder
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	global.Routes.Store(key.String(), serviceID)
	return nil
}

// 删除路由
func DeleteRoute(routeGroupID, routePath, routeMethod string) error {
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
	var key strings.Builder
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	global.Routes.Delete(key.String())
	return nil
}

// 匹配路由，返回集群ID和匹配结果的HTTP状态码
func matchRoute(req *http.Request) (hostname string, service global.ServiceType, status int) {
	var (
		routeGroupID string
		routePath    = req.URL.Path
		routeMethod  = req.Method
		matchResult  bool
		serviceID    string
	)

	// 匹配主机
	hostname, routeGroupID, matchResult = matchHost(req.Host)
	if !matchResult {
		status = http.StatusServiceUnavailable
		return
	}

	if routePath == "" {
		routePath = "/"
	}
	var key strings.Builder
	var exist bool
	// 先尝试直接匹配路由
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	if v, e := global.Routes.Load(key.String()); e {
		service, exist = matchService(v.(string))
		if exist {
			return
		}
	}

	// --------------- 以下用于路径和方法的模糊匹配，并且判断要返回哪种http状态码 -----------------------
	// 先尝试匹配精确路径
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	paths := map[string]string{}
	global.Routes.Range(func(k, v interface{}) bool {
		p := k.(string)
		if strings.HasPrefix(p, key.String()) {
			paths[p] = v.(string)
		}
		return true
	})
	// 如果没有匹配到精确路径，开始尝试匹配模糊路径
	if len(paths) == 0 {
		pathLastChar := routePath[len(routePath)-1]
		// 如果路径的最后一个字符是/
		if pathLastChar != 47 {
			pos := strings.LastIndex(routePath, path.Base(routePath))
			routePath = routePath[:pos]
		}
		routePath = routePath + "*"
		key.Reset()
		key.WriteString(routeGroupID)
		key.WriteString("/")
		key.WriteString(routePath)
		key.WriteString("/")
		// 尝试模糊匹配
		global.Routes.Range(func(k, v interface{}) bool {
			p := k.(string)
			if strings.HasPrefix(p, key.String()) {
				paths[p] = v.(string)
			}
			return true
		})
	}
	// 如果精确和模糊路径都没匹配到，返回404错误
	if len(paths) == 0 {
		status = http.StatusNotFound
		return
	}

	// 尝试匹配精确方法
	for k, v := range paths {
		key.Reset()
		key.WriteString(routeGroupID)
		key.WriteString("/@")
		key.WriteString(routePath)
		key.WriteString("/")
		key.WriteString(routeMethod)
		if strings.HasPrefix(k, key.String()) {
			serviceID = v
			break
		}
	}
	// 尝试匹配ANY方法
	if serviceID == "" {
		for k, v := range paths {
			key.Reset()
			key.WriteString(routeGroupID)
			key.WriteString("/")
			key.WriteString(routePath)
			key.WriteString("/ANY")
			if strings.HasPrefix(k, key.String()) {
				serviceID = v
				break
			}
		}
	}
	if serviceID == "" {
		status = http.StatusMethodNotAllowed
		return
	}

	// 获得服务
	exist = false
	service, exist = matchService(serviceID)
	if !exist {
		status = http.StatusNotImplemented
		return
	}
	return
}
