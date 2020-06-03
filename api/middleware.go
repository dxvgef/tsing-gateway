package api

import (
	"encoding/json"

	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

type Middleware struct{}

func (self *Middleware) Put(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		mw   []global.ModuleConfig
	)
	if err = json.Unmarshal(global.StrToBytes(ctx.Post("middleware")), &mw); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	for k := range mw {
		if _, err = middleware.Build(mw[k].Name, mw[k].Config, true); err != nil {
			resp["error"] = err.Error()
			return JSON(ctx, 400, &resp)
		}
	}
	if err = global.Storage.PutHostMiddleware(ctx.Post("middleware")); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
