package etcd

import (
	"errors"
	"path"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
)

// 从etcd key里解析路由信息
func parseRouteGroup(key []byte) (routeGroupID, routePath, routeMethod string, err error) {
	keyStr := global.TrimPrefix(key, "/routes/")
	pos := strings.Index(keyStr, "/")
	if pos == -1 {
		err = errors.New("路由解析失败")
		return
	}
	routeGroupID = keyStr[:pos]
	if routeGroupID == "" {
		err = errors.New("路由组ID失败")
		return
	}
	routePath = strings.TrimLeft(keyStr, routeGroupID)
	routeMethod = path.Base(routePath)
	if routeMethod != "GET" &&
		routeMethod != "POST" &&
		routeMethod != "PUT" &&
		routeMethod != "DELETE" &&
		routeMethod != "OPTIONS" &&
		routeMethod != "HEAD" &&
		routeMethod != "TRACE" &&
		routeMethod != "PATCH" &&
		routeMethod != "CONNECT" &&
		routeMethod != "*" {
		err = errors.New("路由方法解析失败")
	}
	routePath = strings.TrimRight(routePath, "/"+routeMethod)
	return
}
