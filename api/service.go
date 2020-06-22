package api

import (
	"encoding/json"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

type Service struct{}

func (self *Service) Add(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			id             string
			middleware     string
			discover       string
			staticEndpoint string
			loadBalance    string
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
		filter.El(&req.staticEndpoint, filter.FromString(ctx.Post("static_endpoint"), "static_endpoint")),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}
	if _, exists := global.Services.Load(req.id); exists {
		resp["error"] = "服务ID已存在"
		return JSON(ctx, 400, &resp)
	}
	if req.staticEndpoint == "" && req.discover == "" {
		resp["error"] = "static_endpoint和discover参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}
	if req.discover == "" && req.loadBalance == "" {
		resp["error"] = "discover和load_balance参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}
	if req.staticEndpoint != "" {
		req.discover = ""
		req.loadBalance = ""
	}

	if req.discover != "" {
		if err = service.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
			resp["error"] = "探测器配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &service.Middleware); err != nil {
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	service.ID = req.id
	service.StaticEndpoint = req.staticEndpoint

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("JSON编码失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
func (self *Service) Put(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		req  struct {
			id             string
			middleware     string
			discover       string
			staticEndpoint string
			loadBalance    string
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.PathParams.Value("id"), "id").Required().Base64RawURLDecode()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
		filter.El(&req.staticEndpoint, filter.FromString(ctx.Post("static_endpoint"), "static_endpoint")),
		filter.El(&req.loadBalance, filter.FromString(ctx.Post("load_balance"), "load_balance")),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	if req.staticEndpoint == "" && req.discover == "" {
		resp["error"] = "static_endpoint和discover参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}
	if req.discover == "" && req.loadBalance == "" {
		resp["error"] = "discover和load_balance参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}
	if req.staticEndpoint != "" {
		req.discover = ""
		req.loadBalance = ""
	}

	if req.discover != "" {
		if err = service.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
			resp["error"] = "探测器配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &service.Middleware); err != nil {
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	service.ID = req.id
	service.StaticEndpoint = req.staticEndpoint

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Msg("JSON编码失败")
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}

func (self *Service) Delete(ctx *tsing.Context) error {
	var (
		err  error
		resp = make(map[string]string)
		id   string
	)
	if id, err = global.DecodeKey(ctx.PathParams.Value("id")); err != nil {
		return Status(ctx, 404)
	}
	if _, exist := global.Services.Load(id); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageService(ctx.PathParams.Value("id"))
	if err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
