package api

import (
	"github.com/dxvgef/tsing"
)

type Host struct{}

func (self *Host) Put(ctx *tsing.Context) error {
	resp := make(map[string]string)
	hostname := ctx.Post("hostname")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if err := sa.PutHost(hostname, ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Del(ctx *tsing.Context) error {
	resp := make(map[string]string)
	hostname := ctx.PathParams.Value("hostname")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if err := sa.DelHost(hostname); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
