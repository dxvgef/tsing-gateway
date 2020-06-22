package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
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
				log.Err(err).Caller().Str("hostname", hostname).Str("middleware name", mw[k].GetName()).Msg("执行主机中间件时出错")
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
			log.Error().Str("service id", service.ID).Caller().Msg("服务中间件类型断言失败")
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
				log.Err(err).Caller().Str("service id", service.ID).Str("middleware name", mw[k].GetName()).Msg("执行服务中间件时出错")
				return
			}
			if !next {
				return
			}
		}
	}

	// 获得端点
	var endpointURL *url.URL
	endpointURL, err = getEndpoint(service)
	if err != nil {
		log.Err(err).Caller().Msg("获得端点失败")
		resp.WriteHeader(status)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(status))); err != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 转发请求到端点
	p := httputil.NewSingleHostReverseProxy(endpointURL)
	req.URL.Host = endpointURL.Host
	req.URL.Scheme = endpointURL.Scheme
	p.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		log.Debug().Err(err).Caller().Msg("向端点发起请求出错")
	}
	p.ServeHTTP(resp, req)
}

func getEndpoint(service global.ServiceType) (endpoint *url.URL, err error) {
	if service.StaticEndpoint != "" {
		endpoint, err = url.Parse(service.StaticEndpoint)
		if err != nil {
			log.Err(err).Caller().Msg("解析服务的EndPoint属性失败")
			return nil, err
		}
		return endpoint, nil
	}

	var (
		dc   global.DiscoverType
		ip   string
		port uint16
	)
	dc, err = discover.Build(service.Discover.Name, service.Discover.Config)
	if err != nil {
		log.Err(err).Caller().Msg("构建探测器失败")
		return nil, err
	}
	ip, port, err = dc.Fetch(service.ID)
	if err != nil {
		log.Err(err).Caller().Msg("使用探测器获取节点失败")
		return nil, err
	}
	endpoint, err = url.Parse(ip + ":" + strconv.Itoa(int(port)))
	if err != nil {
		log.Err(err).Caller().Msg("解析探测器返回的节点失败")
		return nil, err
	}
	return endpoint, nil
}
