package api

import (
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"local/global"
	"local/proxy"
)

type Data struct{}

func (self *Data) OutputJSON(ctx *tsing.Context) error {
	data, err := proxy.OutputJSON()
	if err != nil {
		log.Err(err).Caller().Send()
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

func (*Data) LoadAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := loadAll(); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
func (*Data) SaveAll(ctx *tsing.Context) error {
	resp := make(map[string]string)
	if err := saveAll(); err != nil {
		log.Err(err).Caller().Send()
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
