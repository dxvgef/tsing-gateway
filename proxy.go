package main

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/middleware"
)

// 代理引擎
type Proxy struct {
	middleware      []middleware.Middleware                 // 全局中间件
	hosts           map[string]string                       // [主机名]路由组ID
	routeGroups     map[string]map[string]map[string]string // [路由组ID][路径][方法]上游ID
	upstreams       map[string]Upstream                     // [上游ID]上游信息
	hostsUpdated    bool                                    // 主机列表有更新
	routeUpdated    bool                                    // 路由列表有更新
	UpstreamUpdated bool                                    // 上游列表有更新
}

// 获得代理引擎的实例
func newProxy() *Proxy {
	var proxy Proxy
	proxy.hosts = make(map[string]string)
	proxy.routeGroups = make(map[string]map[string]map[string]string)
	proxy.upstreams = make(map[string]Upstream)
	return &proxy
}

// 实现http.handler接口，同时也是下游的请求入口
func (p *Proxy) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var (
		next bool
		err  error
		// endpointURL *url.URL
	)
	upstream, status := p.matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		if _, respErr := resp.Write(strToBytes(http.StatusText(status))); respErr != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 执行全局中间件
	log.Debug().Int("全局中间件数量", len(p.middleware)).Send()
	for k := range p.middleware {
		next, err = p.middleware[k].Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}
		if !next {
			return
		}
	}

	// 执行上游中注册的中间件
	for k := range upstream.Middleware {
		next, err = upstream.Middleware[k].Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}
		if !next {
			return
		}
	}

	// 以下是反向代理的请求逻辑，暂时用200状态码替代
	resp.WriteHeader(http.StatusOK)
	if _, err := resp.Write(strToBytes(http.StatusText(http.StatusOK))); err != nil {
		log.Error().Msg(err.Error())
	}

	// req.Header.Set("X-Forwarded-Host", req.Host)
	// req.Header.Set("X-Power-By", "Tsing Gateway")

	// endpointURL, err = url.Parse(upstream.Endpoints[0].Addr)
	// proxy := httputil.NewSingleHostReverseProxy(endpointURL)
	// req.URL.Host = endpointURL.Host
	// req.URL.Scheme = endpointURL.Scheme
	// req.Host = endpointURL.Host

	// 这里使用的servHTTP是一个使用新协程的非阻塞处理方式
	// resp.Header().Set("X-Power-By", "Tsing Gateway")
	// p.ServeHTTP(resp, req)
}
