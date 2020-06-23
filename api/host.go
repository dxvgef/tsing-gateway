package api

import (
	"strings"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

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
	if _, exists := global.Hosts.Load(hostname); exists {
		resp["error"] = "主机名已存在"
		return JSON(ctx, 400, &resp)
	}
	var host global.HostType
	if err := host.UnmarshalJSON(global.StrToBytes(config)); err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = "config参数解析失败"
		return JSON(ctx, 400, &resp)
	}
	if host.RouteGroupID == "" {
		resp["error"] = "config.route_group_id参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	hostname = global.EncodeKey(strings.ToLower(hostname))
	if err := global.Storage.SaveHost(hostname, config); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Put(ctx *tsing.Context) error {
	var resp = make(map[string]string)
	hostname := ctx.PathParams.Value("hostname")
	config := ctx.Post("config")
	if hostname == "" {
		resp["error"] = "hostname参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	var host global.HostType
	if err := host.UnmarshalJSON(global.StrToBytes(config)); err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = "config参数解析失败"
		return JSON(ctx, 400, &resp)
	}
	if host.RouteGroupID == "" {
		resp["error"] = "config.route_group_id参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if err := global.Storage.SaveHost(hostname, config); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Host) Delete(ctx *tsing.Context) error {
	var (
		hostname string
		resp     = make(map[string]string)
	)
	hostname = ctx.PathParams.Value("hostname")
	if hostname == "" {
		return Status(ctx, 404)
	}
	hostnamePlainText, err := global.DecodeKey(hostname)
	if err != nil {
		// 由于数据来自客户端，因此不记录日志
		return Status(ctx, 500)
	}
	if _, exist := global.Hosts.Load(hostnamePlainText); !exist {
		return Status(ctx, 404)
	}
	if err := global.Storage.DeleteStorageHost(hostname); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
