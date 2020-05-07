package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	var err error

	setDefaultLogger()

	if err = loadConfigFile(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	// reset default logger with local configuration file
	if err = setLogger(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = setEtcdCli(); err != nil {
		log.Fatal().Caller().Msg(err.Error())
	}

	proxy := newProxy()
	proxy.start()
}
