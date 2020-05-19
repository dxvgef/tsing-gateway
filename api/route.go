package api

import (
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
		filter.El(&req.groupID, filter.FromString(ctx.Post("group_id")).Required().UnescapeURLPath()),
		filter.El(&req.path, filter.FromString(ctx.Post("path")).Required()),
		filter.El(&req.method, filter.FromString(ctx.Post("method")).Required().EnumString(global.Methods)),
		filter.El(&req.upstreamID, filter.FromString(ctx.Post("upstream_id")).Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if err = sa.PutRoute(req.groupID, req.path, req.method, req.upstreamID); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Update(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var req struct {
		groupID    string
		path       string
		method     string
		upstreamID string
	}
	err := filter.MSet(
		filter.El(&req.groupID, filter.FromString(ctx.Post("group_id")).Required().UnescapeURLPath()),
		filter.El(&req.path, filter.FromString(ctx.Post("path")).Required()),
		filter.El(&req.method, filter.FromString(ctx.Post("method")).Required().EnumString(global.Methods)),
		filter.El(&req.upstreamID, filter.FromString(ctx.Post("upstream_id")).Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if err = sa.PutRoute(req.groupID, req.path, req.method, req.upstreamID); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Route) Delete(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var req struct {
		groupID string
		path    string
		method  string
	}
	err := filter.MSet(
		filter.El(&req.groupID, filter.FromString(ctx.PathParams.Value("groupID")).Required().UnescapeURLPath()),
		filter.El(&req.path, filter.FromString(ctx.PathParams.Value("path")).Required().UnescapeURLPath()),
		filter.El(&req.method, filter.FromString(ctx.PathParams.Value("method")).Required().EnumString(global.Methods)),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	err = sa.DelRoute(req.groupID, req.path, req.method)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
