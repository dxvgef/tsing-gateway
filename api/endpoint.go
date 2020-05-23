package api

import (
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

type Endpoint struct{}

func (self *Endpoint) Put(ctx *tsing.Context) error {
	var (
		err      error
		hostname string
		resp     = make(map[string]string)
	)
	hostname, err = global.DecodeKey(ctx.PathParams.Value("hostname"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err = global.Storage.PutHost(hostname, ctx.Post("upstream_id")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Endpoint) Delete(ctx *tsing.Context) error {
	var (
		err      error
		hostname string
		resp     = make(map[string]string)
	)
	hostname, err = global.DecodeKey(ctx.PathParams.Value("hostname"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err := global.Storage.DelHost(hostname); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
