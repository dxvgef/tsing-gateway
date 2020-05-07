package main

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// 代理引擎
type Proxy struct {
	hosts           map[string]string                       // [主机名]路由组ID
	routeGroups     map[string]map[string]map[string]string // [路由组ID][路径][方法]上游ID
	upstreams       map[string]Upstream                     // [上游ID]上游信息
	hostsUpdated    bool                                    // 主机列表有更新
	routeUpdated    bool                                    // 路由列表有更新
	UpstreamUpdated bool                                    // 上游列表有更新
}

// 获得新的代理引擎
func New() *Proxy {
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
	)
	upstream, status := p.matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		if _, respErr := resp.Write(strToBytes(http.StatusText(status))); respErr != nil {
			log.Error().Msg(err.Error())
		}
		return
	}
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
}
