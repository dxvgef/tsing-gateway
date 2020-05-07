package main

import (
	"go.etcd.io/etcd/clientv3"
)

var etcdCli *clientv3.Client

func setEtcdCli() (err error) {
	etcdCli, err = clientv3.New(clientv3.Config{
		Endpoints:            localConfig.ETCD.Endpoints,
		DialTimeout:          localConfig.ETCD.DialTimeout,
		Username:             localConfig.ETCD.Username,
		Password:             localConfig.ETCD.Password,
		AutoSyncInterval:     localConfig.ETCD.AutoSyncInterval,
		DialKeepAliveTime:    localConfig.ETCD.DialKeepAliveTime,
		DialKeepAliveTimeout: localConfig.ETCD.DialKeepAliveTimeout,
		MaxCallSendMsgSize:   localConfig.ETCD.MaxCallSendMsgSize,
		MaxCallRecvMsgSize:   localConfig.ETCD.MaxCallRecvMsgSize,
		RejectOldCluster:     localConfig.ETCD.RejectOldCluster,
		PermitWithoutStream:  localConfig.ETCD.PermitWithoutStream,
	})
	return
}
