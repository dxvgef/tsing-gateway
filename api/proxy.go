package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

type Proxy struct{}

func (self *Proxy) OutputJSON(ctx *tsing.Context) error {
	data, err := proxy.OutputJSON()
	if err != nil {
		ctx.ResponseWriter.WriteHeader(500)
		return err
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.ResponseWriter.WriteHeader(200)
	if _, err = ctx.ResponseWriter.Write(data); err != nil {
		log.Err(err).Caller().Send()
		return err
	}
	return nil
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
