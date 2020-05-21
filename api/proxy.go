package api

import (
	"encoding/json"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

type Proxy struct {
	Middleware string                                  `json:"middleware"`
	Hosts      map[string]string                       `json:"hosts"`
	Routes     map[string]map[string]map[string]string `json:"routes"`
	Upstreams  map[string]global.UpstreamType          `json:"upstreams"`
}

func (self *Proxy) OutputAll(ctx *tsing.Context) error {
	self.Hosts = global.Hosts
	self.Routes = global.Routes

	if global.GlobalMiddleware != nil && len(global.GlobalMiddleware) > 0 {
		mw, err := json.Marshal(&global.GlobalMiddleware)
		if err != nil {
			return err
		}
		self.Middleware = global.BytesToStr(mw)
	}

	self.Upstreams = global.Upstreams
	err := JSON(ctx, 200, self)

	self.Hosts = nil
	self.Routes = nil
	self.Upstreams = nil
	self.Middleware = ""
	log.Debug().Caller().Interface("globalMiddleware", global.GlobalMiddleware).Send()
	log.Debug().Caller().Interface("upstreamMiddleware", global.UpstreamMiddleware).Send()
	return err
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
