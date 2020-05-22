package api

import (
	"encoding/base64"

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
		filter.El(&req.method, filter.FromString(ctx.Post("method"), "method").Required().EnumString(global.Methods)),
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
	if err = global.Storage.PutRoute(req.groupID, req.path, req.method, req.upstreamID); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Put(ctx *tsing.Context) error {
	var (
		err          error
		resp         = make(map[string]string)
		key          []byte
		routeGroupID string
		routePath    string
		routeMethod  string
	)
	key, err = base64.RawURLEncoding.DecodeString(ctx.PathParams.Value("key"))
	if err != nil {
		return Status(ctx, 404)
	}
	routeGroupID, routePath, routeMethod, err = global.ParseRoute(global.BytesToStr(key), "")
	if err != nil {
		return Status(ctx, 404)
	}
	if err = global.Storage.PutRoute(routeGroupID, routePath, routeMethod, ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Delete(ctx *tsing.Context) error {
	var (
		resp         = make(map[string]string)
		err          error
		key          []byte
		routeGroupID string
		routePath    string
		routeMethod  string
	)
	key, err = base64.RawURLEncoding.DecodeString(ctx.PathParams.Value("key"))
	if err != nil {
		return Status(ctx, 404)
	}
	routeGroupID, routePath, routeMethod, err = global.ParseRoute(global.BytesToStr(key), "")
	if err != nil {
		return Status(ctx, 404)
	}
	err = global.Storage.DelRoute(routeGroupID, routePath, routeMethod)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
