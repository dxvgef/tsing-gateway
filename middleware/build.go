package middleware

import (
	"errors"

	"github.com/rs/zerolog/log"

	"local/global"
	"local/middleware/auto_response"
	"local/middleware/jwt_proxy"
	"local/middleware/set_header"
	"local/middleware/url_rewrite"
)

// 构建多个中间件实例
// key为中间件名称，value为中间件的参数json字符串
func Build(name, config string, checkExist bool) (global.MiddlewareType, error) {
	switch name {
	case "auto_response":
		f, err := auto_response.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		if checkExist {
			return nil, nil
		}
		return f, nil
	case "set_header":
		if checkExist {
			return nil, nil
		}
		f, err := set_header.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		if checkExist {
			return nil, nil
		}
		return f, nil
	case "url_rewrite":
		if checkExist {
			return nil, nil
		}
		f, err := url_rewrite.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		if checkExist {
			return nil, nil
		}
		return f, nil
	case "jwt_proxy":
		if checkExist {
			return nil, nil
		}
		f, err := jwt_proxy.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		if checkExist {
			return nil, nil
		}
		return f, nil
	}
	return nil, errors.New("不支持的中间件名称 " + name)
}
