package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/storage"
)

func main() {
	var (
		configFile  string
		err         error
		httpServer  *http.Server
		httpsServer *http.Server
		sa          storage.Storage
	)

	// 设置默认logger
	setDefaultLogger()

	// 加载配置文件
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err = global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// 初始化代理引擎实例
	var proxyEngine proxy.Engine
	proxyEngine.Hosts = make(map[string]string)
	proxyEngine.Routes = make(map[string]map[string]map[string]string)
	proxyEngine.Upstreams = make(map[string]proxy.Upstream)

	// 生成唯一ID
	proxyEngine.ID = global.GetIDInt64()
	if proxyEngine.ID == 0 {
		log.Fatal().Caller().Msg("无法自动生成ID标识")
		return
	}

	// 根据配置构建存储器
	sa, err = storage.Build(&proxyEngine, global.Config.Storage.Name, global.Config.Storage.Config)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	log.Debug().Interface("storage", sa).Send()
	// 从存储器中加载所有数据
	if err = sa.LoadAll(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	log.Debug().Interface("proxy", proxyEngine).Send()
	log.Info().Msg("成功从" + global.Config.Storage.Name + "加载数据")

	// 启动HTTP代理
	if global.Config.Proxy.HTTP.Port > 0 {
		httpServer = &http.Server{
			Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTP.Port), 10),
			Handler:           &proxyEngine,
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

	// 启动HTTPS代理
	if global.Config.Proxy.HTTPS.Port > 0 {
		httpsServer = &http.Server{
			Addr:              global.Config.Proxy.IP + ":" + strconv.FormatUint(uint64(global.Config.Proxy.HTTPS.Port), 10),
			Handler:           &proxyEngine,
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

	// 监听存储中的数据变更
	go func() {
		log.Info().Msg("开始监听数据变更")
		if err = sa.Watch(); err != nil {
			log.Fatal().Msg(err.Error())
			return
		}
	}()

	// 阻塞并等待退出超时
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

	// // 启动API服务
	// if global.Config.API.On {
	// 	// 启动grpc服务
	// 	go startApiGrpcServer()
	// 	// 启用grpc gateway服务
	// 	go startApiGrpcGatewayServer()
	// }

	log.Info().Msg("进程已退出")
}

// // 初始化数据，目前仅开发调试用途
// func initData(e *proxy.Engine, st storage.Storage) (err error) {
// 	var (
// 		upstream   proxy.Upstream
// 		routeGroup proxy.RouteGroup
// 	)
// 	upstream.ID = "testUpstream"
// 	upstream.Middleware = append(upstream.Middleware, proxy.Configurator{
// 		Name:   "favicon",
// 		Config: `{"status": 204}`,
// 	})
// 	upstream.Discover.Name = "coredns_etcd"
// 	upstream.Discover.Config = `{"host":"test.uam.local"}`
// 	// 设置上游
// 	err = e.NewUpstream(upstream, false)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
//
// 	// 设置上游
// 	upstream = proxy.Upstream{}
// 	upstream.ID = "test2Upstream"
// 	upstream.Discover.Name = "coredns_etcd"
// 	upstream.Discover.Config = `{"host":"test2.uam.local"}`
// 	err = e.NewUpstream(upstream, false)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
//
// 	// 设置路由组
// 	routeGroup, err = e.SetRouteGroup("testGroup", false)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
// 	// 设置路由
// 	err = routeGroup.SetRoute("/test", "get", "testUpstream", false)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
// 	// 设置主机
// 	err = e.SetHost("127.0.0.1", routeGroup.ID, false)
// 	if err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
//
// 	// 序列化成json
// 	log.Debug().Interface("配置", e).Send()
//
// 	// 将所有数据保存到存储器
// 	if err = st.SaveAll(); err != nil {
// 		log.Err(err).Caller().Send()
// 		return
// 	}
//
// 	return nil
// }
