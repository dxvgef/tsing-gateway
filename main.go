package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
)

func main() {
	var err error

	// 加载本地配置
	if err = loadConfigFile(); err != nil {
		panic(err.Error())
	}

	// 设置logger
	if err = setLogger(); err != nil {
		panic(err.Error())
	}

	// 启动服务
	start()
}

func start() {
	var httpServer *http.Server
	var httpsServer *http.Server
	var endpoints []Endpoint
	var err error
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10080",
		Weight: 100,
	})
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10082",
		Weight: 100,
	})
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10084",
		Weight: 100,
	})

	proxy := New()
	// 添加上游及端点
	if err = proxy.setUpstream(Upstream{
		ID:        "userLogin",
		Endpoints: endpoints,
	}); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "userRegister",
		Endpoints: endpoints,
	}); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "user",
		Endpoints: endpoints,
	}); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "root",
		Endpoints: endpoints,
	}); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	// 添加路由组
	routeGroup, err := proxy.newRouteGroup("uam_v1_routes")
	if err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	// 在路由组内写入路由规则
	if err = routeGroup.setRoute("/user/login", "GET", "userLogin"); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = routeGroup.setRoute("/user/register", "GET", "userRegister"); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = routeGroup.setRoute("/user/*", "GET", "user"); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	if err = routeGroup.setRoute("/", "GET", "root"); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}
	// 添加主机
	if err = proxy.setHost("127.0.0.1", "uam_v1_routes"); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if LocalConfig.Listener.HTTPPort > 0 {
		httpServer = &http.Server{
			Addr:              LocalConfig.Listener.IP + ":" + strconv.Itoa(LocalConfig.Listener.HTTPPort),
			Handler:           proxy,
			ReadTimeout:       LocalConfig.Listener.ReadTimeout,       // 读取超时
			WriteTimeout:      LocalConfig.Listener.WriteTimeout,      // 响应超时
			IdleTimeout:       LocalConfig.Listener.IdleTimeout,       // 连接空闲超时
			ReadHeaderTimeout: LocalConfig.Listener.ReadHeaderTimeout, // header读取超时
		}
		go func() {
			log.Info().Msg("启动 HTTP 代理服务 :8000")
			if err = httpServer.ListenAndServe(); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTP 代理服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
			}
		}()
	}

	if LocalConfig.Listener.HTTPSPort > 0 {
		httpsServer = &http.Server{
			Addr:              LocalConfig.Listener.IP + ":" + strconv.Itoa(LocalConfig.Listener.HTTPSPort),
			Handler:           proxy,
			ReadTimeout:       LocalConfig.Listener.ReadTimeout,       // 读取超时
			WriteTimeout:      LocalConfig.Listener.WriteTimeout,      // 响应超时
			IdleTimeout:       LocalConfig.Listener.IdleTimeout,       // 连接空闲超时
			ReadHeaderTimeout: LocalConfig.Listener.ReadHeaderTimeout, // header读取超时
		}
		go func() {
			log.Info().Msg("启动 HTTPS 代理服务 :8443")
			if LocalConfig.Listener.HTTP2 {
				log.Info().Msg("启用 HTTP/2 支持在 HTTPS 代理服务")
				if err = http2.ConfigureServer(httpsServer, &http2.Server{}); err != nil {
					log.Fatal().Caller().Msg(err.Error())
				}
			}
			if err = httpsServer.ListenAndServeTLS("server.cert", "server.key"); err != nil {
				if err == http.ErrServerClosed {
					log.Info().Msg("HTTPS 代理服务已关闭")
					return
				}
				log.Fatal().Caller().Msg(err.Error())
			}
		}()
	}

	// 退出进程时等待
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// 定义退出超时
	ctx, cancel := context.WithTimeout(context.Background(), LocalConfig.Listener.QuitWaitTimeout)
	defer cancel()

	// 退出http服务
	if LocalConfig.Listener.HTTPPort > 0 {
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
	// 退出https服务
	if LocalConfig.Listener.HTTPSPort > 0 {
		if err := httpsServer.Shutdown(ctx); err != nil {
			log.Fatal().Caller().Msg(err.Error())
		}
	}
}
