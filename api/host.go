package api

import (
	"encoding/base64"

	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
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
		err      error
		hostname []byte
		resp     = make(map[string]string)
	)
	hostname, err = base64.URLEncoding.DecodeString(ctx.PathParams.Value("hostname"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err = sa.PutHost(global.BytesToStr(hostname), ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Delete(ctx *tsing.Context) error {
	var (
		err      error
		hostname []byte
		resp     = make(map[string]string)
	)
	hostname, err = base64.URLEncoding.DecodeString(ctx.PathParams.Value("hostname"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err := sa.DelHost(global.BytesToStr(hostname)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
