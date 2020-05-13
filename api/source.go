package api

import (
	"context"

	"github.com/dxvgef/filter"
	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/source"
)

// 数据源
type SourceHandler struct {
	UnimplementedAPIServer
}

func (*SourceHandler) SetSource(_ context.Context, req *Source) (*Null, error) {
	log.Debug().Interface("req", req).Send()
	return &Null{}, nil
}

// 加载所有数据
func (self *SourceHandler) LoadAll(ctx *tsing.Context) error {
	var (
		err        error
		dataSource source.Source
		// 请求参数
		req struct {
			name       string // 数据源名称
			config     string // 数据源配置(JSON字符串)
			returnData bool   // 返回加载后的数据
		}
	)

	// 响应参数
	resp := map[string]string{
		"error": "",
	}

	// 接收并验证请求参数
	if err = filter.MSet(
		filter.El(&req.name, filter.FromString(ctx.Post("name")).
			Required()),
		filter.El(&req.config, filter.FromString(ctx.Post("config")).
			Required().IsJSON()),
		filter.El(&req.returnData, filter.FromString(ctx.Post("return_data")).
			IsBool()),
	); err != nil {
		resp["error"] = err.Error()
		return JSON(ctx, 400, &resp)
	}

	// 构建数据源实例
	if dataSource, err = source.Build(proxyEngine, req.name, req.config); err != nil {
		log.Err(err).Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	// 加载所有数据
	if err = dataSource.LoadAll(); err != nil {
		log.Err(err).Send()
		resp["error"] = err.Error()
		return JSON(ctx, 500, &resp)
	}

	// 输出所有配置给客户端
	if req.returnData {
		if dataJSON, err := proxyEngine.MarshalJSON(); err != nil {
			return JSONBytes(ctx, 500, dataJSON)
		}
	}
	return nil
}
