package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"

	"local/discover"
	"local/global"
)

// 代理引擎
type Engine struct{}

// 实现http.Handler接口的方法
// 下游请求入口
func (*Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var (
		next bool
		err  error
	)
	// 匹配到服务
	hostname, service, status := matchRoute(req)
	if status > 0 {
		resp.WriteHeader(status)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(status))); err != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 执行主机中间件
	hostMW, hostMWExist := global.HostMiddleware.Load(hostname)
	if hostMWExist {
		mw, ok := hostMW.([]global.MiddlewareType)
		if !ok {
			log.Error().Str("hostname", hostname).Caller().Msg("类型断言失败")
			return
		}
		for k := range mw {
			next = false
			if mw[k] == nil {
				continue
			}
			// 执行中间件逻辑
			next, err = mw[k].Action(resp, req)
			if err != nil {
				log.Err(err).Caller().Str("hostname", hostname).Str("middleware name", mw[k].GetName()).Send()
				return
			}
			if !next {
				return
			}
		}
	}

	// 执行服务中间件
	serviceMW, serviceMWExist := global.ServicesMiddleware.Load(service.ID)
	if serviceMWExist {
		mw, ok := serviceMW.([]global.MiddlewareType)
		if !ok {
			log.Error().Str("service id", service.ID).Caller().Msg("类型断言失败")
			return
		}
		for k := range mw {
			next = false
			if mw[k] == nil {
				continue
			}
			// 执行中间件逻辑
			next, err = mw[k].Action(resp, req)
			if err != nil {
				log.Err(err).Caller().Str("service id", service.ID).Str("middleware name", mw[k].GetName()).Send()
				return
			}
			if !next {
				return
			}
		}
	}

	// 发送数据
	send(service, req, resp)
}

// 获取端点
func getEndpoint(service global.ServiceType) (endpoint *url.URL, err error) {
	if service.StaticEndpoint != "" {
		endpoint, err = url.Parse(service.StaticEndpoint)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		return endpoint, nil
	}

	var (
		dc global.DiscoverType
	)
	dc, err = discover.Build(service.Discover.Name, service.Discover.Config)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	endpoint, err = dc.Fetch(service.ID)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return
}

// 发送数据到端点
func send(service global.ServiceType, req *http.Request, resp http.ResponseWriter) {
	// 获得端点
	var (
		endpointURL *url.URL
		err         error
	)
	endpointURL, err = getEndpoint(service)
	if err != nil {
		log.Err(err).Caller().Send()
		resp.WriteHeader(500)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(500))); err != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 转发请求到端点
	p := httputil.NewSingleHostReverseProxy(endpointURL)
	req.URL.Host = endpointURL.Host
	req.URL.Scheme = endpointURL.Scheme
	p.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		log.Err(err).Caller().Send()
		// if service.Retry > retry && totalTime < global.Config.Proxy.WriteTimeout {
		// 	log.Error().Uint8("retry", retry+1).Caller().Msg("向端点发起重试请求")
		// 	time.Sleep(time.Duration(service.RetryInterval) * time.Millisecond)
		// 	send(service, req, resp, retry+1)
		// 	return
		// }
		resp.WriteHeader(500)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))); err != nil {
			log.Err(err).Caller().Send()
		}
	}
	p.ServeHTTP(resp, req)
}
