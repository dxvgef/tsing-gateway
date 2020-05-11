package main

import (
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

func main() {
	var err error
	setDefaultLogger()

	if err = global.LoadConfigFile(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	// reset default logger with local configuration file
	if err = setLogger(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = global.SetEtcdCli(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	p := NewProxy()
	err = p.LoadConfigFromEtcd()
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}
	return
	// err = p.LoadConfigFromJSON(`{"hosts":{"127.0.0.1":"testGroup"},"route_groups":{"testGroup":{"/user/login":{"GET":"testUpstream"}}},"upstreams":{"testUpstream":{"id":"testUpstream","middleware":[{"name":"favicon","config":"{\"status\": 204}"}],"explorer":{"name":"coredns_etcd","config":"{\"host\":\"test.uam.local\"}"}}}}`)
	// if err != nil {
	// 	log.Fatal().Caller().Msg(err.Error())
	// }
	p.Start()
}
