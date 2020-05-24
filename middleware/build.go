package middleware

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware/favicon"
	"github.com/dxvgef/tsing-gateway/middleware/set_header"
	"github.com/dxvgef/tsing-gateway/middleware/url_rewrite"
)

// 构建多个中间件实例
// key为中间件名称，value为中间件的参数json字符串
func Build(name, config string, test bool) (global.MiddlewareType, error) {
	switch name {
	case "favicon":
		f, err := favicon.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		if test {
			return nil, nil
		}
		return f, nil
	case "set_header":
		if test {
			return nil, nil
		}
		f, err := set_header.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		if test {
			return nil, nil
		}
		return f, nil
	case "url_rewrite":
		if test {
			return nil, nil
		}
		f, err := url_rewrite.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		if test {
			return nil, nil
		}
		return f, nil
	}
	return nil, errors.New("中间件不存在")
}
