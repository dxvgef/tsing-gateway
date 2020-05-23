package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

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
		err         error
		upstream    global.UpstreamType
		status      int
		endpointURL *url.URL
	)
	// 通过路由匹配到上游
	upstream, status = matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		// nolint
		resp.Write(global.StrToBytes(http.StatusText(status)))
		return
	}

	// 执行全局中间件
	for k := range global.GlobalMiddleware {
		next, err := global.GlobalMiddleware[k].Action(resp, req)
		if err != nil {
			log.Err(err).Caller().Str("name", global.GlobalMiddleware[k].GetName()).Msg("执行全局中间件时出错")
			return
		}
		if !next {
			return
		}
	}

	// 执行上游中间件
	for k := range global.UpstreamMiddleware[upstream.ID] {
		// 执行中间件逻辑
		next, err := global.UpstreamMiddleware[upstream.ID][k].Action(resp, req)
		if err != nil {
			log.Err(err).Caller().Str("upstream id", upstream.ID).Str("middleware name", global.UpstreamMiddleware[upstream.ID][k].GetName()).Msg("执行上游中间件时出错")
			return
		}
		if !next {
			return
		}
	}

	// 如果有静态端点
	if upstream.StaticEndpoint != "" {
		endpointURL, err = url.Parse(upstream.StaticEndpoint)
		if err != nil {
			log.Error().Caller().Str("upstream id", upstream.ID).Str("static endpoint", upstream.StaticEndpoint).Msg("静态端点地址不正确")
			resp.WriteHeader(http.StatusBadGateway)
			// nolint
			resp.Write(global.StrToBytes(http.StatusText(http.StatusBadGateway)))
			return
		}
	} else {
		// 构建探测器
		_, err = discover.Build(upstream.Discover.Name, upstream.Discover.Config)
		if err != nil {
			log.Error().Caller().Str("name", upstream.Discover.Name).Str("config", upstream.Discover.Config).Err(err).Msg("构建探测器时出错")
			resp.WriteHeader(http.StatusInternalServerError)
			_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
			return
		}
	}

	// 转发请求到端点
	p := httputil.NewSingleHostReverseProxy(endpointURL)
	req.URL.Host = endpointURL.Host
	req.URL.Scheme = endpointURL.Scheme
	p.ServeHTTP(resp, req)
}
