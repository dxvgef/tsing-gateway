package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/engine"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/source"
)

func main() {
	setDefaultLogger()

	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	err := global.LoadConfigFile(configFile)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	if err = setLogger(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	e := engine.NewEngine()

	dataSrc, err := source.Build(e, global.Config.Source.Name, global.Config.Source.Config)
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	if err = dataSrc.LoadAll(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}
	log.Debug().Interface("加载配置", e)

	e.Start()
}
