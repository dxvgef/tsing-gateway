package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/storage"
)

func main() {
	var (
		configFile string
		err        error
		sa         storage.Storage
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

	// 获得一个代理引擎实例
	proxyEngine := New()

	// 根据配置构建存储器
	sa, err = storage.Build(proxyEngine, global.Config.Storage.Name, global.Config.Storage.Config)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	// 从存储器中加载所有数据
	if err = sa.LoadAll(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// // 启动API服务
	// if global.Config.API.On {
	// 	// 启动grpc服务
	// 	go startApiGrpcServer()
	// 	// 启用grpc gateway服务
	// 	go startApiGrpcGatewayServer()
	// }

	// 启动代理引擎
	start(proxyEngine)
}

// 初始化数据，目前仅开发调试用途
func initData(e *proxy.Engine, st storage.Storage) (err error) {
	var (
		upstream   proxy.Upstream
		routeGroup proxy.RouteGroup
	)
	upstream.ID = "testUpstream"
	upstream.Middleware = append(upstream.Middleware, proxy.Configurator{
		Name:   "favicon",
		Config: `{"status": 204}`,
	})
	upstream.Discover.Name = "coredns_etcd"
	upstream.Discover.Config = `{"host":"test.uam.local"}`
	// 设置上游
	err = e.NewUpstream(upstream, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 设置上游
	upstream = proxy.Upstream{}
	upstream.ID = "test2Upstream"
	upstream.Discover.Name = "coredns_etcd"
	upstream.Discover.Config = `{"host":"test2.uam.local"}`
	err = e.NewUpstream(upstream, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 设置路由组
	routeGroup, err = e.SetRouteGroup("testGroup", false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 设置路由
	err = routeGroup.SetRoute("/test", "get", "testUpstream", false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}
	// 设置主机
	err = e.SetHost("127.0.0.1", routeGroup.ID, false)
	if err != nil {
		log.Err(err).Caller().Send()
		return
	}

	// 序列化成json
	log.Debug().Interface("配置", e).Send()

	// 将所有数据保存到存储器
	if err = st.SaveAll(); err != nil {
		log.Err(err).Caller().Send()
		return
	}

	return nil
}
