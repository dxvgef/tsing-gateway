package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/engine"
	"github.com/dxvgef/tsing-gateway/global"
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

	// 获得一个引擎实例
	e := engine.NewEngine()

	// 构建数据源实例
	dataSource, err := source.Build(e, global.Config.Source.Name, global.Config.Source.Config)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	// 加载所有数据
	if err = dataSource.LoadAll(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	log.Debug().Interface("数据", e)

	// 启动引擎
	e.Start()
}

// 初始化数据，目前仅开发调试用途
func initData(e *engine.Engine, dataSource source.Source) (err error) {
	var (
		upstream   engine.Upstream
		routeGroup engine.RouteGroup
	)
	upstream.ID = "testUpstream"
	upstream.Middleware = append(upstream.Middleware, engine.Configurator{
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
	upstream = engine.Upstream{}
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
