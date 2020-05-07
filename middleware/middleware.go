package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/middleware/favicon"
	"github.com/dxvgef/tsing-gateway/middleware/header"
	"github.com/dxvgef/tsing-gateway/middleware/health"
)

// 定义中间件接口
type Middleware interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
}

// 构建多个中间件实例
// key为中间件名称，value为中间件的参数json字符串
func Build(mw map[string]string) (result []Middleware) {
	for name, config := range mw {
		switch name {
		case "favicon":
			f, err := favicon.New(config)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return
			}
			result = append(result, f)
		case "health":
			f, err := health.New(config)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return
			}
			result = append(result, f)
		case "header":
			f, err := header.New(config)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return
			}
			result = append(result, f)
		}
	}
	return
}
