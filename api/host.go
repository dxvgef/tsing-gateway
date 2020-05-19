package api

import (
	"net/url"

	"github.com/dxvgef/tsing"
)

type Host struct{}

func (self *Host) Add(ctx *tsing.Context) error {
	var resp = make(map[string]string)
	hostname := ctx.Post("hostname")
	upstreamID := ctx.Post("upstream_id")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if upstreamID == "" {
		resp["error"] = "upstream_id参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if _, exists := proxyEngine.Hosts[hostname]; exists {
		resp["error"] = "主机名已存在"
		return JSON(ctx, 400, &resp)
	}
	if err := sa.PutHost(hostname, ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Put(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
	)
	hostname, exists := ctx.PathParams.Get("hostname")
	if !exists {
		return Status(ctx, 404)
	}
	hostname, err = url.PathUnescape(hostname)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if err = sa.PutHost(hostname, ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
	)
	hostname, exists := ctx.PathParams.Get("hostname")
	if !exists {
		return Status(ctx, 404)
	}
	hostname, err = url.PathUnescape(hostname)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if err := sa.DelHost(hostname); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
