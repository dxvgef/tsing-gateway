package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/yaml.v2"

	"github.com/dxvgef/tsing-gateway/global"
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

	p := NewProxy()
	p.Start()
}

// 加载配置文件
func loadConfigFile() error {
	var configPath string
	flag.StringVar(&configPath, "c", "./config.yml", "配置文件路径")
	flag.Parse()
	file, err := os.Open(filepath.Clean(configPath))
	if err != nil {
		return err
	}
	err = yaml.NewDecoder(file).Decode(&global.LocalConfig)
	if err != nil {
		return err
	}
	return nil
}

func setEtcdCli() (err error) {
	global.EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            global.LocalConfig.ETCD.Endpoints,
		DialTimeout:          global.LocalConfig.ETCD.DialTimeout,
		Username:             global.LocalConfig.ETCD.Username,
		Password:             global.LocalConfig.ETCD.Password,
		AutoSyncInterval:     global.LocalConfig.ETCD.AutoSyncInterval,
		DialKeepAliveTime:    global.LocalConfig.ETCD.DialKeepAliveTime,
		DialKeepAliveTimeout: global.LocalConfig.ETCD.DialKeepAliveTimeout,
		MaxCallSendMsgSize:   global.LocalConfig.ETCD.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:   global.LocalConfig.ETCD.MaxCallRecvMsgSize,
		RejectOldCluster:     global.LocalConfig.ETCD.RejectOldCluster,
		PermitWithoutStream:  global.LocalConfig.ETCD.PermitWithoutStream,
	})
	return
}
