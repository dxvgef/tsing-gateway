package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/middleware/favicon"
	"github.com/dxvgef/tsing-gateway/middleware/health_check"
)

// 定义中间件接口
type Middleware interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
}

// 获得中间件实例
// key为中间件名称，value为中间件的参数json字符串
func GetInst(filters map[string]string) (result []Middleware) {
	for name, config := range filters {
		switch name {
		case "favicon":
			f, err := favicon.Inst(config)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return
			}
			result = append(result, f)
		case "health_check":
			f, err := health_check.Inst(config)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return
			}
			result = append(result, f)
		}
	}
	return
}
