package global

import (
	"go.etcd.io/etcd/clientv3"
)

var EtcdCli *clientv3.Client

func SetEtcdCli() (err error) {
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            Config.Etcd.Endpoints,
		DialTimeout:          Config.Etcd.DialTimeout,
		Username:             Config.Etcd.Username,
		Password:             Config.Etcd.Password,
		AutoSyncInterval:     Config.Etcd.AutoSyncInterval,
		DialKeepAliveTime:    Config.Etcd.DialKeepAliveTime,
		DialKeepAliveTimeout: Config.Etcd.DialKeepAliveTimeout,
		MaxCallSendMsgSize:   int(Config.Etcd.MaxCallSendMsgSize),
		MaxCallRecvMsgSize:   int(Config.Etcd.MaxCallRecvMsgSize),
		RejectOldCluster:     Config.Etcd.RejectOldCluster,
		PermitWithoutStream:  Config.Etcd.PermitWithoutStream,
	})
	return
}
