package proxy

import (
	"errors"
	"net/http"
	"strings"

	"local/global"

	"github.com/rs/zerolog/log"
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
		routeGroupID, serviceID string
		routePath               = req.URL.Path
		routeMethod             = req.Method
		matchResult, exist      bool
		key                     strings.Builder
	)
	if routePath == "" {
		routePath = "/"
	}

	// -------------------------------------- 匹配主机 -----------------------------------------------
	hostname, routeGroupID, matchResult = matchHost(req.Host)
	if !matchResult {
		status = http.StatusServiceUnavailable
		return
	}

	// -------------------------------------- 匹配服务 -----------------------------------------------
	// 尝试精确匹配路径和方法
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	if k, v := global.Routes.Load(key.String()); v {
		service, exist = matchService(k.(string))
		if exist {
			return
		}
	}

	// 尝试只精确匹配路径
	key.Reset()
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	var (
		reqGroupID, reqPath, reqMethod string
		err                            error
	)
	// 从请求中解析出信息
	reqGroupID, reqPath, reqMethod, err = global.ParseRouteFromKey(key.String()+routeMethod, "")
	if err != nil {
		status = http.StatusInternalServerError
		// 此处error会由客户端请求触发，因此不记录日志
		return
	}
	// 根据精准路径匹配服务
	serviceID, status, err = matchPath(false, reqGroupID, reqPath, reqMethod)
	if err != nil {
		log.Err(err)
	}
	if serviceID == "" {
		// 根据通配路径匹配服务
		serviceID, status, err = matchPath(true, reqGroupID, reqPath, reqMethod)
	}
	if err != nil {
		status = http.StatusInternalServerError
		// 此处error会由客户端请求触发，因此不记录日志
		return
	}
	if status != 0 {
		return
	}

	// 匹配服务
	exist = false
	service, exist = matchService(serviceID)
	if !exist {
		status = http.StatusNotImplemented
		return
	}
	return
}

// 匹配路径
func matchPath(isWildcard bool, reqGroupID, reqPath, reqMethod string) (serviceID string, status int, err error) {
	var (
		keyGroupID, keyPath, keyMethod string
	)
	global.Routes.Range(func(k, v interface{}) bool {
		keyStr := k.(string)
		// 从key中解析出信息
		keyGroupID, keyPath, keyMethod, err = global.ParseRouteFromKey(keyStr, "")
		if err != nil {
			log.Err(err).Caller().Send()
			return false
		}
		// 对比路由组
		if keyGroupID != reqGroupID {
			return true
		}

		// 匹配精准路径
		if !isWildcard && reqPath == keyPath {
			if keyMethod == "ANY" {
				serviceID = v.(string)
				return false
			}
			if keyMethod == reqMethod {
				serviceID = v.(string)
				return false
			}
			status = http.StatusMethodNotAllowed
			return false
		}
		// 匹配通配路径
		if isWildcard && strings.HasPrefix(reqPath, keyPath[:len(keyPath)-1]) {
			if keyMethod == "ANY" {
				serviceID = v.(string)
				return false
			}
			if keyMethod == reqMethod {
				serviceID = v.(string)
				return false
			}
			status = http.StatusMethodNotAllowed
			return false
		}
		return true
	})
	return
}
