package main

import (
	"go.etcd.io/etcd/clientv3"
)

var EtcdClient *clientv3.Client

func setEtcdClient() (err error) {
	EtcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.50.212:2379"},
		DialTimeout: LocalConfig.ETCD.DialTimeout,
	})
	return
}
