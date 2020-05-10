package middleware

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/middleware/favicon"
	"github.com/dxvgef/tsing-gateway/middleware/health"
	"github.com/dxvgef/tsing-gateway/middleware/set_header"
)

// 定义中间件接口
type Middleware interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
}

// 构建多个中间件实例
// key为中间件名称，value为中间件的参数json字符串
func Build(name, config string) (Middleware, error) {
	switch name {
	case "favicon":
		f, err := favicon.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		return f, nil
	case "health":
		f, err := health.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		return f, nil
	case "set_header":
		f, err := set_header.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("middleware not found by name")
}
