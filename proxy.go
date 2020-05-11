package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/dxvgef/tsing-gateway/explorer"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

// 代理引擎的数据
type Proxy struct {
	id              int64
	Middleware      []Configurator                          `json:"middleware,omitempty"`   // 全局中间件
	Hosts           map[string]string                       `json:"hosts,omitempty"`        // [hostname]routeGroupID
	RouteGroups     map[string]map[string]map[string]string `json:"route_groups,omitempty"` // [routeGroupID][reqPath][reqMethod]upstreamID
	Upstreams       map[string]Upstream                     `json:"upstreams,omitempty"`    // [upstreamID]Host
	hostsUpdated    bool
	routeUpdated    bool
	upstreamUpdated bool
}

// 初始化一个新的代理引擎
func NewProxy() *Proxy {
	var p Proxy
	p.Hosts = make(map[string]string)
	p.RouteGroups = make(map[string]map[string]map[string]string)
	p.Upstreams = make(map[string]Upstream)
	return &p
}

// 实现http.Handler接口的方法
// 下游请求入口
func (p *Proxy) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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
	e, err := explorer.Build(upstream.Explorer.Name, upstream.Explorer.Config)
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

// 启动代理服务
func (p *Proxy) Start() {
	var httpProxy *http.Server
	var httpsProxy *http.Server
	var err error

	p.id = global.GetIDInt64()
	if p.id == 0 {
		log.Fatal().Caller().Msg("无法自动生成ID标识")
		return
	}

	// 启动HTTP代理
	if global.Config.HTTP.Port > 0 {
		httpProxy = &http.Server{
			Addr:              global.Config.IP + ":" + strconv.FormatUint(uint64(global.Config.HTTP.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.HTTP.ReadTimeout,
			WriteTimeout:      global.Config.HTTP.WriteTimeout,
			IdleTimeout:       global.Config.HTTP.IdleTimeout,
			ReadHeaderTimeout: global.Config.HTTP.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("启动HTTP代理 " + httpProxy.Addr)
			if err = httpProxy.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP代理已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// start HTTPS proxy
	if global.Config.HTTPS.Port > 0 {
		httpsProxy = &http.Server{
			Addr:              global.Config.IP + ":" + strconv.FormatUint(uint64(global.Config.HTTPS.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.HTTPS.ReadTimeout,
			WriteTimeout:      global.Config.HTTPS.WriteTimeout,
			IdleTimeout:       global.Config.HTTPS.IdleTimeout,
			ReadHeaderTimeout: global.Config.HTTPS.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("启动HTTPS代理 " + httpsProxy.Addr)
			if global.Config.HTTPS.HTTP2 {
				log.Info().Msg("启用HTTP2代理支持")
				if err = http2.ConfigureServer(httpsProxy, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}
			if err = httpsProxy.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS代理已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// 等待退出超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), global.Config.QuitWaitTimeout)
	defer cancel()

	// 关闭HTTP服务
	if global.Config.HTTP.Port > 0 {
		if err := httpProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
	// 关闭HTTPS服务
	if global.Config.HTTPS.Port > 0 {
		if err := httpsProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
}
