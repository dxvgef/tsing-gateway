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
			// retry          uint8
			// retryInterval  uint16
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.Post("id"), "id").Required()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
		filter.El(&req.staticEndpoint, filter.FromString(ctx.Post("static_endpoint"), "static_endpoint")),
		// filter.El(&req.retry, filter.FromString(ctx.Post("retry"), "retry").IsDigit().MinInteger(0).MaxInteger(math.MaxUint8)),
		// filter.El(&req.retryInterval, filter.FromString(ctx.Post("retry_interval"), "retry_interval").IsDigit().MinInteger(0).MaxInteger(math.MaxUint16)),
	); err != nil {
		// 由于数据来自客户端，因此不记录日志
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
	if req.staticEndpoint != "" {
		req.discover = ""
	}

	if req.discover != "" {
		if err = service.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
			// 由于数据来自客户端，因此不记录日志
			resp["error"] = "探测器配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &service.Middleware); err != nil {
			// 由于数据来自客户端，因此不记录日志
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	service.ID = req.id
	service.StaticEndpoint = req.staticEndpoint
	// service.Retry = req.retry
	// service.RetryInterval = req.retryInterval

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
		log.Err(err).Caller().Send()
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
			// retry          uint8
			// retryInterval  uint16
		}
		service      global.ServiceType
		serviceBytes []byte
	)
	if err = filter.MSet(
		filter.El(&req.id, filter.FromString(ctx.PathParams.Value("id"), "id").Required().Base64RawURLDecode()),
		filter.El(&req.discover, filter.FromString(ctx.Post("discover"), "discover").IsJSON()),
		filter.El(&req.middleware, filter.FromString(ctx.Post("middleware"), "middleware").IsJSON()),
		filter.El(&req.staticEndpoint, filter.FromString(ctx.Post("static_endpoint"), "static_endpoint")),
		// filter.El(&req.retry, filter.FromString(ctx.Post("retry"), "retry").IsDigit().MinInteger(0).MaxInteger(math.MaxUint8)),
		// filter.El(&req.retryInterval, filter.FromString(ctx.Post("retry_interval"), "retry_interval").IsDigit().MinInteger(0).MaxInteger(math.MaxUint16)),
	); err != nil {
		// 由于数据来自客户端，因此不记录日志
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	if req.staticEndpoint == "" && req.discover == "" {
		resp["error"] = "static_endpoint和discover参数不能同时为空"
		return JSON(ctx, 400, &resp)
	}
	if req.staticEndpoint != "" {
		req.discover = ""
	}

	if req.discover != "" {
		if err = service.Discover.UnmarshalJSON(global.StrToBytes(req.discover)); err != nil {
			// 由于数据来自客户端，因此不记录日志
			resp["error"] = "探测器配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	if req.middleware != "" {
		if err = json.Unmarshal(global.StrToBytes(req.middleware), &service.Middleware); err != nil {
			// 由于数据来自客户端，因此不记录日志
			resp["error"] = "中间件配置解析失败"
			return JSON(ctx, 400, &resp)
		}
	}

	service.ID = req.id
	service.StaticEndpoint = req.staticEndpoint
	// service.Retry = req.retry
	// service.RetryInterval = req.retryInterval
	// log.Debug().Uint8("retry", service.Retry).Caller().Send()

	if serviceBytes, err = service.MarshalJSON(); err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	if err = global.Storage.SaveService(req.id, global.BytesToStr(serviceBytes)); err != nil {
		log.Err(err).Caller().Send()
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
		// 由于数据来自客户端，因此不记录日志
		return Status(ctx, 404)
	}
	if _, exist := global.Services.Load(id); !exist {
		return Status(ctx, 404)
	}
	err = global.Storage.DeleteStorageService(ctx.PathParams.Value("id"))
	if err != nil {
		log.Err(err).Caller().Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}
	return Status(ctx, 204)
}
