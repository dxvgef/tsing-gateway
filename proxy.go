package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"

	"github.com/dxvgef/tsing-gateway/middleware"
)

// proxy engine
type Proxy struct {
	middleware      []middleware.Middleware                 // global middleware
	hosts           map[string]string                       // [hostname]routeGroupID
	routeGroups     map[string]map[string]map[string]string // [routeGroupID][reqPath][reqMethod]upstreamID
	upstreams       map[string]Upstream                     // [upstreamID]Upstream
	hostsUpdated    bool                                    // hosts map changed
	routeUpdated    bool                                    // routeGroups map changed
	UpstreamUpdated bool                                    // upstreams map changed
}

// get instance of proxy engine
func newProxy() *Proxy {
	var proxy Proxy
	proxy.hosts = make(map[string]string)
	proxy.routeGroups = make(map[string]map[string]map[string]string)
	proxy.upstreams = make(map[string]Upstream)
	return &proxy
}

// implement http.Handler interface
// downstream request entry
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

	// execute global middleware
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

	// execute upstream middleware
	for k := range upstream.Middleware {
		next, err = upstream.Middleware[k].Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
		}
		if !next {
			return
		}
	}

	// todo 以下是反向代理的请求逻辑，暂时用200状态码替代
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

// start proxy engine
func (p *Proxy) start() {
	var httpProxy *http.Server
	var httpsProxy *http.Server
	var err error

	log.Info().Msg("listen IP: [" + localConfig.Proxy.IP + "]")
	// listen HTTP request
	if localConfig.Proxy.HTTP.Port > 0 {
		httpPort := strconv.Itoa(localConfig.Proxy.HTTP.Port)
		httpProxy = &http.Server{
			Addr:              localConfig.Proxy.IP + ":" + httpPort,
			Handler:           p,
			ReadTimeout:       localConfig.Proxy.ReadTimeout,
			WriteTimeout:      localConfig.Proxy.WriteTimeout,
			IdleTimeout:       localConfig.Proxy.IdleTimeout,
			ReadHeaderTimeout: localConfig.Proxy.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("HTTP proxy port: [" + httpPort + "]")
			if err = httpProxy.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP proxy is down")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
			}
		}()
	}

	// listen HTTPS request
	if localConfig.Proxy.HTTPS.Port > 0 {
		httpsPort := strconv.Itoa(localConfig.Proxy.HTTPS.Port)
		httpsProxy = &http.Server{
			Addr:              localConfig.Proxy.IP + ":" + httpsPort,
			Handler:           p,
			ReadTimeout:       localConfig.Proxy.ReadTimeout,
			WriteTimeout:      localConfig.Proxy.WriteTimeout,
			IdleTimeout:       localConfig.Proxy.IdleTimeout,
			ReadHeaderTimeout: localConfig.Proxy.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Msg("HTTPS proxy port: [" + httpsPort + "]")
			if localConfig.Proxy.HTTPS.HTTP2 {
				log.Info().Msg("HTTP2 support is enabled")
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

	ctx, cancel := context.WithTimeout(context.Background(), localConfig.Proxy.QuitWaitTimeout)
	defer cancel()

	// shutdown the HTTP service
	if localConfig.Proxy.HTTP.Port > 0 {
		if err := httpProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
	// shutdown the HTTPS service
	if localConfig.Proxy.HTTPS.Port > 0 {
		if err := httpsProxy.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
}
