package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

func main() {
	var err error
	setDefaultLogger()

	var configFile string
	flag.StringVar(&configFile, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	if err = global.LoadConfigFile(configFile); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if err = setLogger(); err != nil {
		log.Fatal().Msg(err.Error())
		return
	}

	if err = global.SetEtcdCli(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	p := NewProxy()

	err = p.LoadDataFromEtcd()
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return
	}

	p.Start()
}
