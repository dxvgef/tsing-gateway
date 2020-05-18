package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

func checkHeader(ctx *tsing.Context) error {
	log.Debug().Str("secret", ctx.Request.Header.Get("SECRET")).Msg("检查secfret")
	if ctx.Request.Header.Get("SECRET") != global.Config.API.Secret {
		ctx.Abort()
		return Status(ctx, 401)
	}
	return nil
}
