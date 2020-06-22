package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

func SetService(service global.ServiceType) error {
	if service.ID == "" {
		return errors.New("service ID不能为空")
	}
	global.Services.Store(service.ID, service)

	// 更新中间件
	mwLen := len(service.Middleware)
	if mwLen == 0 {
		global.ServicesMiddleware.Delete(service.ID)
		return nil
	}
	mw := make([]global.MiddlewareType, mwLen)
	for k := range service.Middleware {
		m, err := middleware.Build(service.Middleware[k].Name, service.Middleware[k].Config, false)
		if err != nil {
			return err
		}
		mw = append(mw, m)
	}
	global.ServicesMiddleware.Store(service.ID, mw)
	return nil
}

func DelService(serviceID string) error {
	global.Services.Delete(serviceID)
	global.ServicesMiddleware.Delete(serviceID)
	return nil
}

func matchService(serviceID string) (global.ServiceType, bool) {
	if serviceID == "" {
		return global.ServiceType{}, false
	}
	service, exist := global.Services.Load(serviceID)
	if !exist {
		return global.ServiceType{}, false
	}
	return service.(global.ServiceType), true
}
