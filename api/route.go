package api

import (
	"strings"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

type Route struct{}

func (self *Route) Add(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var req struct {
		groupID    string
		path       string
		method     string
		upstreamID string
	}
	err := filter.MSet(
		filter.El(&req.groupID, filter.FromString(ctx.Post("group_id"), "group_id").Required()),
		filter.El(&req.path, filter.FromString(ctx.Post("path"), "path").Required()),
		filter.El(&req.method, filter.FromString(ctx.Post("method"), "method").Required().ToUpper().EnumString(global.Methods)),
		filter.El(&req.upstreamID, filter.FromString(ctx.Post("upstream_id"), "upstream_id").Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exist := global.Routes[req.groupID][req.path][req.method]; exist {
		resp["error"] = "路由已存在"
		return JSON(ctx, 400, &resp)
	}
	if err = global.Storage.PutRoute(req.groupID, req.path, req.method, req.upstreamID, false); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Update(ctx *tsing.Context) error {
	var (
		err          error
		resp         = make(map[string]string)
		key          strings.Builder
		routeGroupID string
		routePath    string
		routeMethod  string
	)
	key.WriteString("/")
	key.WriteString(ctx.PathParams.Value("groupID"))
	key.WriteString("/")
	key.WriteString(ctx.PathParams.Value("path"))
	key.WriteString("/")
	key.WriteString(ctx.PathParams.Value("method"))
	routeGroupID, routePath, routeMethod, err = global.ParseRoute(key.String(), "")
	if err != nil {
		return Status(ctx, 404)
	}
	if _, exist := global.Routes[routeGroupID][routePath][routeMethod]; !exist {
		return Status(ctx, 404)
	}
	if err = global.Storage.PutRoute(routeGroupID, routePath, routeMethod, ctx.Post("upstream_id"), false); err != nil {
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
	)
	err = filter.MSet(
		filter.El(&routeGroupID, filter.FromString(ctx.PathParams.Value("groupID"), "group_id").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
		filter.El(&routeMethod, filter.FromString(ctx.PathParams.Value("method"), "method").Required().ToUpper().EnumString(global.Methods)),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exist := global.Routes[routeGroupID][routePath][routeMethod]; !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DelRoute(routeGroupID, routePath, routeMethod, false)
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
	)
	err = filter.MSet(
		filter.El(&routeGroupID, filter.FromString(ctx.PathParams.Value("groupID"), "group_id").Required().Base64RawURLDecode()),
		filter.El(&routePath, filter.FromString(ctx.PathParams.Value("path"), "path").Required().Base64RawURLDecode()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exist := global.Routes[routeGroupID][routePath]; !exist {
		return Status(ctx, 404)
	}
	if len(global.Routes[routeGroupID][routePath]) == 0 {
		return Status(ctx, 404)
	}
	err = global.Storage.DelRoute(routeGroupID, routePath, "", false)
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
	)
	routeGroupID, err = global.DecodeKey(routeGroupID)
	if err != nil {
		return Status(ctx, 404)
	}
	if _, exist := global.Routes[routeGroupID]; !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DelRoute(routeGroupID, "", "", false)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
	)
	global.Routes = make(map[string]map[string]map[string]string)
	err = global.Storage.DelRoute("", "", "", false)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
