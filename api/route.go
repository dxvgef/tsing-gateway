package api

import (
	"strings"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

type Route struct{}

func (self *Route) Add(ctx *tsing.Context) error {
	var (
		resp = map[string]string{}
		req  struct {
			groupID    string
			path       string
			method     string
			upstreamID string
		}
		key strings.Builder
	)
	err := filter.MSet(
		filter.El(&req.groupID, filter.FromString(ctx.Post("group_id"), "group_id").Required()),
		filter.El(&req.path, filter.FromString(ctx.Post("path"), "path").Required()),
		filter.El(&req.method, filter.FromString(ctx.Post("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
		filter.El(&req.upstreamID, filter.FromString(ctx.Post("upstream_id"), "upstream_id").Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(req.groupID)
	key.WriteString("/")
	key.WriteString(req.path)
	key.WriteString("/")
	key.WriteString(req.method)
	if _, exist := global.Routes.Load(key.String()); exist {
		resp["error"] = "路由已存在"
		return JSON(ctx, 400, &resp)
	}
	if err = global.Storage.SaveRoute(req.groupID, req.path, req.method, req.upstreamID); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Put(ctx *tsing.Context) error {
	var (
		resp = map[string]string{}
		req  struct {
			groupID    string
			path       string
			method     string
			upstreamID string
		}
		key strings.Builder
	)
	err := filter.MSet(
		filter.El(&req.groupID, filter.FromString(ctx.PathParams.Value("groupID"), "groupID").Required().Base64RawURLDecode()),
		filter.El(&req.path, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
		filter.El(&req.method, filter.FromString(ctx.PathParams.Value("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
		filter.El(&req.upstreamID, filter.FromString(ctx.Post("upstream_id"), "upstream_id").Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(req.groupID)
	key.WriteString("/")
	key.WriteString(req.path)
	key.WriteString("/")
	key.WriteString(req.method)
	if err = global.Storage.SaveRoute(req.groupID, req.path, req.method, req.upstreamID); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeleteMethod(ctx *tsing.Context) error {
	var (
		err          error
		resp         = make(map[string]string)
		routeGroupID string
		routePath    string
		routeMethod  string
		key          strings.Builder
	)
	err = filter.MSet(
		filter.El(&routeGroupID, filter.FromString(ctx.PathParams.Value("groupID"), "group_id").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
		filter.El(&routeMethod, filter.FromString(ctx.PathParams.Value("method"), "method").Required().ToUpper().EnumString(global.HTTPMethods)),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	if _, exist := global.Routes.Load(key.String()); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageRoute(routeGroupID, routePath, routeMethod)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeletePath(ctx *tsing.Context) error {
	var (
		err          error
		resp         = make(map[string]string)
		routeGroupID string
		routePath    string
		key          strings.Builder
		exist        bool
	)
	err = filter.MSet(
		filter.El(&routeGroupID, filter.FromString(ctx.PathParams.Value("groupID"), "group_id").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	key.WriteString(routeGroupID)
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
	err = global.Storage.DeleteStorageRoute(routeGroupID, routePath, "")
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) DeleteGroup(ctx *tsing.Context) error {
	var (
		err          error
		resp         = make(map[string]string)
		routeGroupID = ctx.PathParams.Value("groupID")
		key          strings.Builder
		exist        bool
	)
	routeGroupID, err = global.DecodeKey(routeGroupID)
	if err != nil {
		return Status(ctx, 404)
	}
	key.WriteString(routeGroupID)
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
	err = global.Storage.DeleteStorageRoute(routeGroupID, "", "")
	if err != nil {
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
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
