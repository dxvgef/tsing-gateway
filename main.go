package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	apiEngine "github.com/dxvgef/tsing-gateway/api"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/source"
)

func main() {
	// 设置默认logger
	setDefaultLogger()

	// 加载配置文件
	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err := global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	// 获得一个代理引擎实例
	proxyEngine := proxy.New()

	// 启动api服务
	if global.Config.API.On {
		go apiEngine.Start(proxyEngine)
	}

	// 启动网关引擎
	proxyEngine.Start()
}

// 初始化数据，目前仅开发调试用途
func initData(e *proxy.Engine, dataSource source.Source) (err error) {
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

	// 将所有数据保存到数据源
	if err = dataSource.SaveAll(); err != nil {
		log.Err(err).Caller().Send()
		return
	}

	return nil
}
