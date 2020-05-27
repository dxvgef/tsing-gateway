package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/load_balance"
)

// 代理引擎
type Engine struct{}

// 实现http.Handler接口的方法
// 下游请求入口
func (*Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var (
		next      bool
		err       error
		endpoints []global.EndpointType
	)
	// 通过路由匹配到上游
	upstream, status := matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		// nolint
		resp.Write(global.StrToBytes(http.StatusText(status)))
		return
	}

	// 执行全局中间件
	for k := range global.GlobalMiddleware {
		next = false
		next, err = global.GlobalMiddleware[k].Action(resp, req)
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
		next = false
		// 执行中间件逻辑
		next, err = global.UpstreamMiddleware[upstream.ID][k].Action(resp, req)
		if err != nil {
			log.Err(err).Caller().Str("upstream id", upstream.ID).Str("middleware name", global.UpstreamMiddleware[upstream.ID][k].GetName()).Msg("执行上游中间件时出错")
			return
		}
		if !next {
			return
		}
	}

	// 获得所有端点列表
	endpoints, status = FetchEndpoints(upstream)
	if status > 0 {
		resp.WriteHeader(status)
		// nolint
		resp.Write(global.StrToBytes(http.StatusText(status)))
	}

	// 使用负载均衡算法选取一个
	lb := load_balance.Build("WR")
	for i := range endpoints {
		lb.Put(endpoints[i].Addr, endpoints[i].Weight)
	}
	var (
		endpointAddr string
		endpointURL  *url.URL
	)
	endpointAddr = lb.Next()
	endpointURL, err = url.Parse(endpointAddr)
	if err != nil {
		log.Err(err).Str("addr", endpointAddr).Caller().Msg("解析端点地址失败")
		resp.WriteHeader(http.StatusInternalServerError)
		// nolint
		resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError)))
		return
	}

	log.Debug().Caller().Str("addr", endpointAddr).Msg("目标端点")
	// 转发请求到端点
	p := httputil.NewSingleHostReverseProxy(endpointURL)
	req.URL.Host = endpointURL.Host
	req.URL.Scheme = endpointURL.Scheme
	p.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {

	}
	p.ServeHTTP(resp, req)
}

// 使用上游的探测器获得最新的端点列表
func FetchEndpoints(upstream global.UpstreamType) (endpoints []global.EndpointType, status int) {
	var (
		err error
		dc  global.DiscoverType
	)
	// 如果有静态端点则静态优先
	if upstream.StaticEndpoint != "" {
		endpoints = append(endpoints, global.EndpointType{
			UpstreamID: upstream.ID,
			Addr:       upstream.StaticEndpoint,
			Weight:     0,
		})
		return
	}

	// 构建探测器
	dc, err = discover.Build(upstream.Discover.Name, upstream.Discover.Config)
	if err != nil {
		log.Err(err).Caller().Str("name", upstream.Discover.Name).Str("config", upstream.Discover.Config).Msg("构建探测器时出错")
		status = http.StatusInternalServerError
		return
	}
	// 获得所有站点
	endpoints, err = dc.Fetch(upstream.ID)
	if err != nil {
		log.Err(err).Caller().Str("name", upstream.Discover.Name).Str("config", upstream.Discover.Config).Msg("获得端点列表时出错")
		status = http.StatusServiceUnavailable
	}
	return
}
