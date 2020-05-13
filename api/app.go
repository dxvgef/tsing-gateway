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

	"github.com/dxvgef/tsing-gateway/api/handler"
	"github.com/dxvgef/tsing-gateway/global"
)

var app *tsing.Engine

func Start() {
	var (
		appServer *http.Server
		rootPath  string
		err       error
		appConfig tsing.Config
	)

	appConfig.EventHandler = handler.EventHandler
	appConfig.Recover = true
	appConfig.EventShortPath = true
	appConfig.EventSource = true
	appConfig.EventTrace = true
	appConfig.EventHandlerError = true
	rootPath, err = os.Getwd()
	if err == nil {
		appConfig.RootPath = rootPath
	}
	app = tsing.New(&appConfig)

	// 设置路由
	setRouter()

	// 定义HTTP服务
	appServer = &http.Server{
		Addr:    global.Config.API.IP + ":" + strconv.Itoa(global.Config.API.Port),
		Handler: app,
		// ErrorLog:          global.Logger.StdError,                                                    // 日志记录器
		// ReadTimeout:       time.Duration(global.Config.Service.ReadTimeout) * time.Second,       // 读取超时
		// WriteTimeout:      time.Duration(global.Config.Service.WriteTimeout) * time.Second,      // 响应超时
		// IdleTimeout:       time.Duration(global.Config.Service.IdleTimeout) * time.Second,       // 连接空闲超时
		// ReadHeaderTimeout: time.Duration(global.Config.Service.ReadHeaderTimeout) * time.Second, // http header读取超时
	}
	// 在新协程中启动服务，方便实现退出等待
	go func() {
		log.Info().Msg("启动API服务 " + global.Config.API.Path + " => " + appServer.Addr)
		if err := appServer.ListenAndServe(); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(global.Config.QuitWaitTimeout)*time.Second)
	defer cancel()
	if err := appServer.Shutdown(ctx); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	log.Info().Msg("API服务已退出")
}
