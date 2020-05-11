package global

import (
	"go.etcd.io/etcd/clientv3"
)

var EtcdCli *clientv3.Client

func SetEtcdCli() (err error) {
	EtcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            LocalConfig.Etcd.Endpoints,
		DialTimeout:          LocalConfig.Etcd.DialTimeout,
		Username:             LocalConfig.Etcd.Username,
		Password:             LocalConfig.Etcd.Password,
		AutoSyncInterval:     LocalConfig.Etcd.AutoSyncInterval,
		DialKeepAliveTime:    LocalConfig.Etcd.DialKeepAliveTime,
		DialKeepAliveTimeout: LocalConfig.Etcd.DialKeepAliveTimeout,
		MaxCallSendMsgSize:   LocalConfig.Etcd.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:   LocalConfig.Etcd.MaxCallRecvMsgSize,
		RejectOldCluster:     LocalConfig.Etcd.RejectOldCluster,
		PermitWithoutStream:  LocalConfig.Etcd.PermitWithoutStream,
	})
	return
}
