package api

import (
	"encoding/json"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

type Upstream struct{}

func (self *Upstream) Add(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var (
		req struct {
			id         string
			middleware string
			discover   string
		}
		upstream      proxy.Upstream
		upstreamBytes []byte
	)
	err := filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.Post("id"), "id").Required().UnescapeURLPath()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").Required().IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	if _, exists := proxyEngine.Upstreams[req.id]; exists {
		resp["error"] = "上游ID已存在"
		return JSON(ctx, 400, &resp)
	}

	if err = upstream.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
		resp["error"] = "探测器配置解析失败"
		return JSON(ctx, 400, &resp)
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &upstream.Middleware); err != nil {
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	upstream.ID = req.id

	if upstreamBytes, err = upstream.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("对upstream序列化成JSON字符串失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = sa.PutUpstream(req.id, global.BytesToStr(upstreamBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
func (self *Upstream) Put(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var (
		req struct {
			id         string
			middleware string
			discover   string
		}
		upstream      proxy.Upstream
		upstreamBytes []byte
	)
	err := filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.PathParams.Value("id"), "id").Required().UnescapeURLPath()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").Required().IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
	)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	if err = upstream.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
		resp["error"] = "探测器配置解析失败"
		return JSON(ctx, 400, &resp)
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &upstream.Middleware); err != nil {
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	upstream.ID = req.id

	if upstreamBytes, err = upstream.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("对upstream序列化成JSON字符串失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = sa.PutUpstream(req.id, global.BytesToStr(upstreamBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Upstream) Del(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		id   string
	)
	id, err = filter.FromString(ctx.PathParams.Value("id"), "id").Required().UnescapeURLPath().String()
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	err = sa.DelUpstream(id)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
