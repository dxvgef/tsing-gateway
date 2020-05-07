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
	upstream, status := p.matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		if _, err := resp.Write(StrToBytes(http.StatusText(status))); err != nil {
			log.Error().Msg(err.Error())
		}
		return
	}
	resp.WriteHeader(http.StatusOK)
	if _, err := resp.Write(StrToBytes(http.StatusText(http.StatusOK))); err != nil {
		log.Error().Msg(err.Error())
	}
	log.Debug().Caller().Str("upstream id", upstream.ID).Send()
}
