package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

var apiEngine *tsing.Engine

var proxyEngine *proxy.Engine

func Start(p *proxy.Engine) {
	proxyEngine = p

	var (
		err         error
		httpServer  *http.Server
		httpsServer *http.Server
		rootPath    string
		apiConfig   tsing.Config
	)

	apiConfig.EventHandler = EventHandler
	apiConfig.Recover = true
	apiConfig.EventShortPath = true
	apiConfig.EventSource = true
	apiConfig.EventTrace = true
	apiConfig.EventHandlerError = true
	rootPath, err = os.Getwd()
	if err == nil {
		apiConfig.RootPath = rootPath
	}
	apiEngine = tsing.New(&apiConfig)

	// 设置路由
	setRouter()

	// 启动HTTP服务
	if global.Config.API.HTTP.Port > 0 {
		httpServer = &http.Server{
			Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTP.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.API.ReadTimeout,
			WriteTimeout:      global.Config.API.WriteTimeout,
			IdleTimeout:       global.Config.API.IdleTimeout,
			ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
		}
		go func() {
			if err = httpServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP API服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// 启动HTTPS API服务
	if global.Config.API.HTTPS.Port > 0 {
		httpsServer = &http.Server{
			Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTPS.Port), 10),
			Handler:           p,
			ReadTimeout:       global.Config.API.ReadTimeout,
			WriteTimeout:      global.Config.API.WriteTimeout,
			IdleTimeout:       global.Config.API.IdleTimeout,
			ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
		}
		go func() {
			if global.Config.API.HTTPS.HTTP2 {
				if err = http2.ConfigureServer(httpsServer, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}
			if err = httpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS API服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	var serverStatus strings.Builder
	serverStatus.WriteString("启动API服务")
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

	ctx, cancel := context.WithTimeout(context.Background(), global.Config.API.QuitWaitTimeout)
	defer cancel()

	// 关闭HTTP API服务
	if global.Config.API.HTTP.Port > 0 {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
	// 关闭HTTPS API服务
	if global.Config.API.HTTPS.Port > 0 {
		if err := httpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
}
