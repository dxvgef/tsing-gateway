package api

import (
	"encoding/json"

	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

type Host struct{}

func (self *Host) Add(ctx *tsing.Context) error {
	var resp = make(map[string]string)
	hostname := ctx.Post("hostname")
	config := ctx.Post("config")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if !json.Valid(global.StrToBytes(config)) {
		resp["error"] = "config参数不是有效的JSON字符串"
		return JSON(ctx, 400, &resp)
	}
	if _, exists := global.Hosts.Load(hostname); exists {
		resp["error"] = "主机名已存在"
		return JSON(ctx, 400, &resp)
	}
	if err := global.Storage.SaveHost(hostname, config); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Put(ctx *tsing.Context) error {
	var resp = make(map[string]string)
	hostname := ctx.Post("hostname")
	config := ctx.Post("config")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if !json.Valid(global.StrToBytes(config)) {
		resp["error"] = "config参数不是有效的JSON字符串"
		return JSON(ctx, 400, &resp)
	}
	if err := global.Storage.SaveHost(hostname, config); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Delete(ctx *tsing.Context) error {
	var (
		err      error
		hostname string
		resp     = make(map[string]string)
	)
	hostname, err = global.DecodeKey(ctx.PathParams.Value("hostname"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err := global.Storage.DeleteStorageHost(hostname); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
