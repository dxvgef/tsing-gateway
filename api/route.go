package api

import (
	"strings"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"local/global"
)

type Route struct{}

func (self *Route) Add(ctx *tsing.Context) error {
	var (
		resp = map[string]string{}
		req  struct {
			hostname  string
			path      string
			method    string
			serviceID string
		}
		key strings.Builder
	)
	err := filter.MSet(
		filter.El(&req.hostname, filter.FromString(ctx.Post("hostname"), "hostname").Required()),
		filter.El(&req.path, filter.FromString(ctx.Post("path"), "path").Required()),
		filter.El(&req.method, filter.FromString(ctx.Post("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
		filter.El(&req.serviceID, filter.FromString(ctx.Post("service_id"), "service_id").Required()),
	)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(req.hostname)
	key.WriteString("/")
	key.WriteString(req.path)
	key.WriteString("/")
	key.WriteString(req.method)
	if _, exist := global.Routes.Load(key.String()); exist {
		resp["error"] = "路由已存在"
		return JSON(ctx, 400, &resp)
	}
	if err = global.Storage.SaveRoute(req.hostname, req.path, req.method, req.serviceID); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Put(ctx *tsing.Context) error {
	var (
		resp = map[string]string{}
		req  struct {
			hostname  string
			path      string
			method    string
			serviceID string
		}
		key strings.Builder
	)
	err := filter.MSet(
		filter.El(&req.hostname, filter.FromString(ctx.PathParams.Value("hostname"), "hostname").Required().Base64RawURLDecode()),
		filter.El(&req.path, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
		filter.El(&req.method, filter.FromString(ctx.PathParams.Value("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
		filter.El(&req.serviceID, filter.FromString(ctx.Post("service_id"), "service_id").Required()),
	)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(req.hostname)
	key.WriteString("/")
	key.WriteString(req.path)
	key.WriteString("/")
	key.WriteString(req.method)
	if err = global.Storage.SaveRoute(req.hostname, req.path, req.method, req.serviceID); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeleteMethod(ctx *tsing.Context) error {
	var (
		err           error
		resp          = make(map[string]string)
		routeHostname string
		routePath     string
		routeMethod   string
		key           strings.Builder
	)
	err = filter.MSet(
		filter.El(&routeHostname, filter.FromString(ctx.PathParams.Value("hostname"), "hostname").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
		filter.El(&routeMethod, filter.FromString(ctx.PathParams.Value("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
	)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(routeHostname)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	if _, exist := global.Routes.Load(key.String()); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageRoute(routeHostname, routePath, routeMethod)
	if err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeletePath(ctx *tsing.Context) error {
	var (
		err           error
		resp          = make(map[string]string)
		routeHostname string
		routePath     string
		key           strings.Builder
		exist         bool
	)
	err = filter.MSet(
		filter.El(&routeHostname, filter.FromString(ctx.PathParams.Value("hostname"), "hostname").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
	)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(routeHostname)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	global.Routes.Range(func(k, v interface{}) bool {
		if strings.HasPrefix(k.(string), key.String()) {
			exist = true
			return false
		}
		return true
	})
	if !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageRoute(routeHostname, routePath, "")
	if err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeleteGroup(ctx *tsing.Context) error {
	var (
		err           error
		resp          = make(map[string]string)
		routeHostname = ctx.PathParams.Value("hostname")
		key           strings.Builder
		exist         bool
	)
	routeHostname, err = global.DecodeKey(routeHostname)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		return Status(ctx, 404)
	}
	key.WriteString(routeHostname)
	key.WriteString("/")
	global.Routes.Range(func(k, v interface{}) bool {
		if strings.HasPrefix(k.(string), key.String()) {
			exist = true
			return false
		}
		return true
	})
	if !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageRoute(routeHostname, "", "")
	if err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeleteAll(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
	)
	global.SyncMapClean(&global.Routes)
	err = global.Storage.DeleteStorageRoute("", "", "")
	if err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
