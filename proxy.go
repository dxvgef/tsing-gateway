package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

// proxy engine
type Proxy struct {
	Middleware      []Configurator                          `json:"middleware,omitempty"`   // global Middleware
	Hosts           map[string]string                       `json:"hosts,omitempty"`        // [hostname]routeGroupID
	RouteGroups     map[string]map[string]map[string]string `json:"route_groups,omitempty"` // [routeGroupID][reqPath][reqMethod]upstreamID
	Upstreams       map[string]Upstream                     `json:"upstreams,omitempty"`    // [upstreamID]Host
	hostsUpdated    bool                                    // Hosts map changed
	routeUpdated    bool                                    // RouteGroups map changed
	upstreamUpdated bool                                    // Upstreams map changed
}

// get instance of proxy engine
func NewProxy() *Proxy {
	var p Proxy
	p.Hosts = make(map[string]string)
	p.RouteGroups = make(map[string]map[string]map[string]string)
	p.Upstreams = make(map[string]Upstream)
	return &p
}

// implement http.Handler interface
// downstream request entry
func (p *Proxy) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var (
		next bool
		err  error
		// endpointURL *url.URL
	)
	upstream, status := p.MatchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		if _, respErr := resp.Write(global.StrToBytes(http.StatusText(status))); respErr != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// execute global Middleware
	log.Debug().Int("全局中间件数量", len(p.Middleware)).Send()
	for k := range p.Middleware {
		mw, err := middleware.Build(p.Middleware[k].Name, p.Middleware[k].Config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		next, err = mw.Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// execute upstream Middleware
	for k := range upstream.Middleware {
		mw, err := middleware.Build(upstream.Middleware[k].Name, upstream.Middleware[k].Config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		next, err = mw.Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// todo 以下是反向代理的请求逻辑，暂时用200状态码替代
	resp.WriteHeader(http.StatusOK)
	if _, err := resp.Write(global.StrToBytes(http.StatusText(http.StatusOK))); err != nil {
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

// start proxy engine
func (p *Proxy) Start() {
	var httpProxy *http.Server
	var httpsProxy *http.Server
	var err error

	// start HTTP proxy
	if global.LocalConfig.HTTP.Port > 0 {
		httpProxy = &http.Server{
			Addr:              global.LocalConfig.IP + ":" + strconv.Itoa(global.LocalConfig.HTTP.Port),
			Handler:           p,
			ReadTimeout:       global.LocalConfig.HTTP.ReadTimeout,
			WriteTimeout:      global.LocalConfig.HTTP.WriteTimeout,
			IdleTimeout:       global.LocalConfig.HTTP.IdleTimeout,
			ReadHeaderTimeout: global.LocalConfig.HTTP.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("start HTTP proxy " + httpProxy.Addr)
			if err = httpProxy.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP proxy is down")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
			}
		}()
	}

	// start HTTPS proxy
	if global.LocalConfig.HTTPS.Port > 0 {
		httpsProxy = &http.Server{
			Addr:              global.LocalConfig.IP + ":" + strconv.Itoa(global.LocalConfig.HTTPS.Port),
			Handler:           p,
			ReadTimeout:       global.LocalConfig.HTTPS.ReadTimeout,
			WriteTimeout:      global.LocalConfig.HTTPS.WriteTimeout,
			IdleTimeout:       global.LocalConfig.HTTPS.IdleTimeout,
			ReadHeaderTimeout: global.LocalConfig.HTTPS.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("start HTTPS proxy " + httpsProxy.Addr)
			if global.LocalConfig.HTTPS.HTTP2 {
				log.Info().Msg("HTTP2 proxy support is enabled")
				if err = http2.ConfigureServer(httpsProxy, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
				}
			}
			if err = httpsProxy.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS proxy is down")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
			}
		}()
	}

	// timeout for waiting to exit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), global.LocalConfig.QuitWaitTimeout)
	defer cancel()

	// shutdown the HTTP service
	if global.LocalConfig.HTTP.Port > 0 {
		if err := httpProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
	// shutdown the HTTPS service
	if global.LocalConfig.HTTPS.Port > 0 {
		if err := httpsProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
}
