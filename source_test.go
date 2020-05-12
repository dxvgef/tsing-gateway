package main

import (
	"flag"
	"testing"

	"github.com/dxvgef/tsing-gateway/engine"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/source"
)

func TestRoute(t *testing.T) {
	var (
		err        error
		configFile string
		dataSource source.Source
	)
	flag.StringVar(&configFile, "c", "./config.local.yml", "配置文件路径")
	flag.Parse()
	if err = global.LoadConfigFile(configFile); err != nil {
		t.Error(err.Error())
		return
	}
	if err = setLogger(); err != nil {
		t.Error(err.Error())
		return
	}

	// 获得一个引擎实例
	e := engine.NewEngine()

	// 构建数据源实例
	dataSource, err = source.Build(e, global.Config.Source.Name, global.Config.Source.Config)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// 初始化数据
	if err = initData(e, dataSource); err != nil {
		t.Error(err.Error())
		return
	}
}
