package api

import (
	"encoding/json"
	"net/http"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"
)

// 输出JSON数据给客户端
func JSON(ctx *tsing.Context, status int, data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(dataBytes); err != nil {
		log.Err(err).Caller().Send()
	}
	return nil
}

// 输出JSON数据给客户端
func JSONBytes(ctx *tsing.Context, status int, data []byte) error {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(status)
	if _, err := ctx.ResponseWriter.Write(data); err != nil {
		log.Err(err).Caller().Send()
	}
	return nil
}

// 输出String给客户端
func String(ctx *tsing.Context, status int, data string) error {
	// 设置header信息
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(status)
	if _, err := ctx.ResponseWriter.Write([]byte(data)); err != nil {
		log.Err(err).Caller().Send()
	}
	return nil
}

// 输出HTTP状态码，无返回数据
func Status(ctx *tsing.Context, status int) error {
	ctx.ResponseWriter.WriteHeader(status)
	if _, err := ctx.ResponseWriter.Write([]byte(http.StatusText(status))); err != nil {
		log.Err(err).Caller().Send()
	}
	return nil
}
