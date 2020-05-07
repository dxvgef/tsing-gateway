package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	var err error
	if err = loadConfigFile(); err != nil {
		panic(err.Error())
	}
	if err = setLogger(); err != nil {
		panic(err.Error())
	}
	if err = setEtcdCli(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	proxy := newProxy() // // get instance of proxy engine
	proxy.start()
}
