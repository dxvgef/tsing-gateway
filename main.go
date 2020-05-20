package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/dxvgef/tsing"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"

	"github.com/dxvgef/tsing-gateway/api"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/storage"
)

func main() {
	var (
		configFile       string
		err              error
		proxyEngine      proxy.Engine
		proxyHttpServer  *http.Server
		proxyHttpsServer *http.Server
		apiHttpServer    *http.Server
		apiHttpsServer   *http.Server
		sa               storage.Storage
	)

	// 设置默认logger
	setDefaultLogger()

	// --------------------- 加载配置文件 ----------------------
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err = global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// --------------------- 初始化代理引擎实例 ----------------------
	proxyEngine.Hosts = make(map[string]string)
	proxyEngine.Routes = make(map[string]map[string]map[string]string)
	proxyEngine.Upstreams = make(map[string]proxy.Upstream)

	// 生成唯一ID
	global.ID = global.GetIDInt64()
	if global.ID == 0 {
		log.Fatal().Caller().Msg("无法自动生成ID标识")
		return
	}

	// --------------------- 根据配置构建存储器 ----------------------
	sa, err = storage.Build(&proxyEngine, global.Config.Storage.Name, global.Config.Storage.Config)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	// 从存储器中加载所有数据
	if err = sa.LoadAll(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// 监听存储中的数据变更
	go func() {
		log.Info().Msg("开始监听数据变更")
		if err = sa.Watch(); err != nil {
			log.Fatal().Msg(err.Error())
			return
		}
	}()

	// 启动HTTP代理
	if global.Config.Proxy.HTTP.Port > 0 {
		proxyHttpServer = &http.Server{
			Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTP.Port), 10),
			Handler:           &proxyEngine,
			ReadTimeout:       global.Config.Proxy.ReadTimeout,
			WriteTimeout:      global.Config.Proxy.WriteTimeout,
			IdleTimeout:       global.Config.Proxy.IdleTimeout,
			ReadHeaderTimeout: global.Config.Proxy.ReadHeaderTimeout,
		}
		go func() {
			log.Info().Str("addr", proxyHttpServer.Addr).Msg("网关HTTP服务")
			if err = proxyHttpServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP代理服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// 启动HTTPS代理
	if global.Config.Proxy.HTTPS.Port > 0 {
		go func() {
			proxyHttpsServer = &http.Server{
				Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTPS.Port), 10),
				Handler:           &proxyEngine,
				ReadTimeout:       global.Config.Proxy.ReadTimeout,
				WriteTimeout:      global.Config.Proxy.WriteTimeout,
				IdleTimeout:       global.Config.Proxy.IdleTimeout,
				ReadHeaderTimeout: global.Config.Proxy.ReadHeaderTimeout,
			}
			if global.Config.Proxy.HTTPS.HTTP2 {
				if err = http2.ConfigureServer(proxyHttpsServer, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}
			log.Info().Bool("HTTP2", global.Config.Proxy.HTTPS.HTTP2).Str("addr", proxyHttpsServer.Addr).Msg("网关HTTPS服务")
			if err = proxyHttpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS代理服务已关闭")
					return
				}
				log.Fatal().Caller().Caller().Msg(err.Error())
				return
			}
		}()
	}

	// 启动API服务
	if global.Config.API.HTTP.Port > 0 || global.Config.API.HTTPS.Port > 0 {
		var (
			apiEngineConfig tsing.Config
			rootPath        string
		)
		apiEngineConfig.EventHandler = api.EventHandler
		apiEngineConfig.Recover = true
		apiEngineConfig.EventShortPath = true
		apiEngineConfig.EventSource = true
		apiEngineConfig.EventTrace = true
		apiEngineConfig.EventHandlerError = true
		rootPath, err = os.Getwd()
		if err == nil {
			apiEngineConfig.RootPath = rootPath
		}
		apiEngine := tsing.New(&apiEngineConfig)
		// 设置存储器
		api.SetStorage(sa)
		// 设置代理引擎
		api.SetProxyEngine(&proxyEngine)
		// 设置路由
		api.SetRouter(apiEngine)
		// 启动api http服务
		if global.Config.API.HTTP.Port > 0 {
			go func() {
				apiHttpServer = &http.Server{
					Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTP.Port), 10),
					Handler:           apiEngine,
					ReadTimeout:       global.Config.API.ReadTimeout,
					WriteTimeout:      global.Config.API.WriteTimeout,
					IdleTimeout:       global.Config.API.IdleTimeout,
					ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
				}
				log.Info().Str("addr", apiHttpServer.Addr).Msg("API HTTP服务")
				if err = apiHttpServer.ListenAndServe(); err != nil {
					if err == http.ErrServerClosed {
						log.Info().Msg("API HTTP服务已关闭")
						return
					}
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}()
		}

		// 启动api https服务
		if global.Config.API.HTTPS.Port > 0 {
			go func() {
				apiHttpsServer = &http.Server{
					Addr:              global.Config.API.IP + ":" + strconv.FormatUint(uint64(global.Config.API.HTTPS.Port), 10),
					Handler:           apiEngine,
					ReadTimeout:       global.Config.API.ReadTimeout,
					WriteTimeout:      global.Config.API.WriteTimeout,
					IdleTimeout:       global.Config.API.IdleTimeout,
					ReadHeaderTimeout: global.Config.API.ReadHeaderTimeout,
				}
				if global.Config.API.HTTPS.HTTP2 {
					if err = http2.ConfigureServer(apiHttpsServer, &http2.Server{}); err != nil {
						log.Fatal().Caller().Msg(err.Error())
						return
					}
				}
				log.Info().Bool("HTTP2", global.Config.API.HTTPS.HTTP2).Str("addr", apiHttpsServer.Addr).Msg("API HTTPS服务")
				if err = apiHttpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
					if err == http.ErrServerClosed {
						log.Info().Msg("API HTTPS服务已关闭")
						return
					}
					log.Fatal().Caller().Msg(err.Error())
					return
				}
			}()
		}
	}

	// 阻塞并等待退出超时
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), global.Config.Proxy.QuitWaitTimeout)
	defer cancel()

	// 关闭网关HTTP服务
	if global.Config.Proxy.HTTP.Port > 0 {
		if err := proxyHttpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
	// 关闭网关HTTPS服务
	if global.Config.Proxy.HTTPS.Port > 0 {
		if err := proxyHttpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}

	// 关闭API HTTP服务
	if global.Config.API.HTTP.Port > 0 {
		if err := apiHttpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}
	// 关闭API HTTPS服务
	if global.Config.API.HTTPS.Port > 0 {
		if err := apiHttpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return
		}
	}

	log.Info().Msg("进程已退出")
}
