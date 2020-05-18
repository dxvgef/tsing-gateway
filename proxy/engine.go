package proxy

import (
	"net/http"
	"strconv"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"

	"github.com/rs/zerolog/log"
)

// 代理引擎
type Engine struct {
	ID              int64                                   `json:"-"`
	Middleware      []global.ModuleConfig                   `json:"middleware,omitempty"` // 全局中间件
	Hosts           map[string]string                       `json:"hosts,omitempty"`      // [hostname]routeGroupID
	Routes          map[string]map[string]map[string]string `json:"routes,omitempty"`     // [routeGroupID][path][method]upstreamID
	Upstreams       map[string]Upstream                     `json:"upstreams,omitempty"`  // [upstreamID]Upstream
	hostsUpdated    bool
	routeUpdated    bool
	upstreamUpdated bool
}

// 实现http.Handler接口的方法
// 下游请求入口
func (p *Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	upstream, status := p.MatchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		// nolint
		_, _ = resp.Write(global.StrToBytes(http.StatusText(status)))
		return
	}

	// 执行全局中间件
	for k := range p.Middleware {
		mw, err := middleware.Build(p.Middleware[k].Name, p.Middleware[k].Config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
			return
		}
		next, err := mw.Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// 执行上游中间件
	for k := range upstream.Middleware {
		mw, err := middleware.Build(upstream.Middleware[k].Name, upstream.Middleware[k].Config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
			return
		}
		next, err := mw.Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// 执行探测器获取端点
	e, err := discover.Build(upstream.Discover.Name, upstream.Discover.Config)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
		return
	}
	ip, port, weight, ttl, err := e.Action()
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
	}
	log.Error().Caller().Str("ip", ip).Int("port", port).Int("weight", weight).Int("ttl", ttl).Send()
	if ip == "" || port == 0 || weight == 0 {
		log.Error().Caller().Str("err", "invalid endpoint").
			Str("ip", ip).Int("port", port).Int("weight", weight).Send()
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
	}

	// todo 以下是反向代理的请求逻辑，暂时用200状态码替代
	respText := `{"ip": "` + ip + `", "port":` + strconv.Itoa(port) + `, "weight":` + strconv.Itoa(weight) + `, "ttl": ` + strconv.Itoa(ttl) + `}`
	resp.WriteHeader(http.StatusOK)
	if _, err := resp.Write(global.StrToBytes(respText)); err != nil {
		log.Error().Msg(err.Error())
	}

	// 这里使用的servHTTP是一个使用新协程的非阻塞处理方式
	// resp.Header().Set("X-Power-By", "Tsing Gateway")
	// p.ServeHTTP(resp, req)
}
