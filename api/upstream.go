package api

import (
	"encoding/base64"
	"encoding/json"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

type Upstream struct{}

func (self *Upstream) Add(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			id         string
			middleware string
			discover   string
		}
		upstream      proxy.Upstream
		upstreamBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").Required().IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
	); err != nil {
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
	var (
		err error
		req struct {
			id         []byte
			middleware string
			discover   string
		}
		resp          = make(map[string]string)
		upstream      proxy.Upstream
		upstreamBytes []byte
	)
	req.id, err = base64.URLEncoding.DecodeString(ctx.PathParams.Value("id"))
	if err != nil {
		return Status(ctx, 404)
	}
	if err = filter.MSet(
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").Required().IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
	); err != nil {
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

	upstream.ID = global.BytesToStr(req.id)

	if upstreamBytes, err = upstream.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("对upstream序列化成JSON字符串失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = sa.PutUpstream(upstream.ID, global.BytesToStr(upstreamBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Upstream) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		id   []byte
	)
	id, err = base64.URLEncoding.DecodeString(ctx.PathParams.Value("id"))
	if err != nil {
		return Status(ctx, 404)
	}
	err = sa.DelUpstream(global.BytesToStr(id))
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
