package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

type Upstream struct{}

func (self *Upstream) Put(ctx *tsing.Context) error {
	resp := make(map[string]string)
	var (
		upstream    proxy.Upstream
		upstreamStr = ctx.Post("upstream")
	)
	err := upstream.UnmarshalJSON(global.StrToBytes(upstreamStr))
	if err != nil {
		log.Err(err).Caller().Msg("解析中间件配置时出错")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	if upstream.ID == "" {
		resp["error"] = "ID参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	if err = sa.PutUpstream(upstream.ID, upstreamStr); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Upstream) Del(ctx *tsing.Context) error {
	resp := make(map[string]string)
	id := ctx.PathParams.Value("id")
	if id == "" {
		resp["error"] = "id参数不能为空"
		return JSON(ctx, 400, &resp)
	}
	err := sa.DelUpstream(id)
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
