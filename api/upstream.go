package api

import (
	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/proxy"
)

type Upstream struct{}

func (self *Upstream) Put(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var upstream proxy.Upstream
	err := filter.MSet(
		filter.El(upstream.ID, filter.FromString(ctx.Post("upstream_id"), "upstream_id").Required()),
		filter.El(upstream.Middleware, filter.FromString(ctx.Post("middleware"), "middleware").Required()),
		filter.El(upstream.Discover, filter.FromString(ctx.Post("discover"), "discover").Required()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	if err = sa.PutUpstream(upstream); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Upstream) Del(ctx *tsing.Context) error {
	resp := make(map[string]string)
	hostname := ctx.PathParams.Value("hostname")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	err := sa.DelUpstream(ctx.PathParams.Value("name"))
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
