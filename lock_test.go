package main

import (
	"testing"

	"github.com/dxvgef/tsing-gateway/global"

	"github.com/coreos/etcd/clientv3/concurrency"
)

func TestLock(t *testing.T) {
	setDefaultLogger()
	global.SetEtcdCli()
	sess, err := concurrency.NewSession(global.EtcdCli)
	if err != nil {
		t.Error(err.Error())
		return
	}
	concurrency.NewMutex(sess, "/tsing-gateway/lock/")
}
