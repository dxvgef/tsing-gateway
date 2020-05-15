package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

// 初始化一个新的代理引擎
func New() *proxy.Engine {
	var p proxy.Engine
	p.Hosts = make(map[string]string)
	p.Routes = make(map[string]map[string]map[string]string)
	p.Upstreams = make(map[string]proxy.Upstream)
	return &p
}

// 启动代理服务
func start(p *proxy.Engine) {
	var (
		err         error
		httpServer  *http.Server
		httpsServer *http.Server
	)

	// 生成唯一ID
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

	// 打印启动信息
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
