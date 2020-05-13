package proxy

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

// 代理引擎
type Engine struct {
	ID              int64                                   `json:"-"`
	Middleware      []Configurator                          `json:"middleware,omitempty"` // 全局中间件
	Hosts           map[string]string                       `json:"hosts,omitempty"`      // [hostname]routeGroupID
	Routes          map[string]map[string]map[string]string `json:"routes,omitempty"`     // [routeGroupID][path][method]upstreamID
	Upstreams       map[string]Upstream                     `json:"upstreams,omitempty"`  // [upstreamID]Upstream
	hostsUpdated    bool
	routeUpdated    bool
	upstreamUpdated bool
}

// 初始化一个新的代理引擎
func New() *Engine {
	var p Engine
	p.Hosts = make(map[string]string)
	p.Routes = make(map[string]map[string]map[string]string)
	p.Upstreams = make(map[string]Upstream)
	return &p
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
func (p *Engine) Start() {
	var httpServer *http.Server
	var httpsServer *http.Server
	var err error

	p.ID = global.GetIDInt64()
	if p.ID == 0 {
		log.Fatal().Caller().Msg("无法自动生成ID标识")
		return
	}

	// 启动HTTP代理
	if global.Config.Proxy.HTTP.Port > 0 {
		httpServer = &http.Server{
			Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTP.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.Proxy.ReadTimeout,
			WriteTimeout:      global.Config.Proxy.WriteTimeout,
			IdleTimeout:       global.Config.Proxy.IdleTimeout,
			ReadHeaderTimeout: global.Config.Proxy.ReadHeaderTimeout,
		}
		go func() {
			if err = httpServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP代理服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// start HTTPS proxy
	if global.Config.Proxy.HTTPS.Port > 0 {
		httpsServer = &http.Server{
			Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTPS.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.Proxy.ReadTimeout,
			WriteTimeout:      global.Config.Proxy.WriteTimeout,
			IdleTimeout:       global.Config.Proxy.IdleTimeout,
			ReadHeaderTimeout: global.Config.Proxy.ReadHeaderTimeout,
		}
		go func() {
			if global.Config.Proxy.HTTPS.HTTP2 {
				if err = http2.ConfigureServer(httpsServer, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}
			if err = httpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS代理服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	var serverStatus strings.Builder
	serverStatus.WriteString("启动代理服务")
	if global.Config.Proxy.HTTP.Port > 0 {
		serverStatus.WriteString(" http://")
		serverStatus.WriteString(httpServer.Addr)
	}
	if global.Config.Proxy.HTTPS.Port > 0 {
		serverStatus.WriteString(" 和 https://")
		serverStatus.WriteString(httpsServer.Addr)
	}
	if global.Config.Proxy.HTTPS.HTTP2 {
		serverStatus.WriteString(" 已启用HTTP2")
	}
	log.Info().Msg(serverStatus.String())

	// 等待退出超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), global.Config.Proxy.QuitWaitTimeout)
	defer cancel()

	// 关闭HTTP服务
	if global.Config.Proxy.HTTP.Port > 0 {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
	// 关闭HTTPS服务
	if global.Config.Proxy.HTTPS.Port > 0 {
		if err := httpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
}
