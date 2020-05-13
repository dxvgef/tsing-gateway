package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

var apiEngine *tsing.Engine

var proxyEngine *proxy.Engine

func Start(p *proxy.Engine) {
	proxyEngine = p

	var (
		err        error
		httpServer *http.Server
		rootPath   string
		apiConfig  tsing.Config
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

	// 定义HTTP服务
	httpServer = &http.Server{
		Addr:    global.Config.API.IP + ":" + strconv.Itoa(global.Config.API.Port),
		Handler: apiEngine,
		// ErrorLog:          global.Logger.StdError,                                                    // 日志记录器
		// ReadTimeout:       time.Duration(global.Config.Service.ReadTimeout) * time.Second,       // 读取超时
		// WriteTimeout:      time.Duration(global.Config.Service.WriteTimeout) * time.Second,      // 响应超时
		// IdleTimeout:       time.Duration(global.Config.Service.IdleTimeout) * time.Second,       // 连接空闲超时
		// ReadHeaderTimeout: time.Duration(global.Config.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
	}
	// 在新协程中启动服务，方便实现退出等待
	go func() {
		log.Info().Msg("启动API服务 " + global.Config.API.Path + " => " + httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info().Msg("API服务已关闭")
				return
			}
			log.Fatal().Msg(err.Error())
		}
	}()

	// 退出进程时等待
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 指定退出超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.Proxy.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	log.Info().Msg("API服务已退出")
}
