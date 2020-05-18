package api

import (
	"github.com/dxvgef/tsing"

	"github.com/dxvgef/tsing-gateway/global"
)

func checkSecret(ctx *tsing.Context) error {
	if ctx.Request.Method == "GET" && ctx.Query("secret") != global.Config.API.Secret {
		ctx.Abort()
		return Status(ctx, 401)
	}
	if ctx.Request.Method != "GET" && ctx.Post("secret") != global.Config.API.Secret {
		ctx.Abort()
		return Status(ctx, 401)
	}
	return nil
}
