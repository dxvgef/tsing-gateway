package api

import (
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

type Proxy struct {
}

func (*Proxy) OutputAll(ctx *tsing.Context) error {
	var proxy global.ProxyStruct
	proxy.Hosts = global.Hosts
	proxy.Routes = global.Routes
	proxy.Middleware = global.Middleware
	proxy.Upstreams = global.Upstreams
	return JSON(ctx, 200, &proxy)
}

func (*Proxy) LoadAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := loadAll(); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
func (*Proxy) SaveAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := saveAll(); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

// 加载所有数据
func loadAll() (err error) {
	return global.Storage.LoadAll()
}

// 保存所有数据
func saveAll() (err error) {
	return global.Storage.SaveAll()
}
