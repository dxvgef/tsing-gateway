package etcd

import (
	"context"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

// 监听变更
func (self *Etcd) Watch() error {
	ch := self.client.Watch(context.Background(), self.KeyPrefix+"/", clientv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			switch event.Type {
			// 更新事件
			case clientv3.EventTypePut:
				if err := self.watchLoadData(event.Kv.Key, event.Kv.Value); err != nil {
					log.Err(err).Caller().Msg("更新本地数据时出错")
				}
			// 删除事件
			case clientv3.EventTypeDelete:
				if err := self.watchDeleteData(event.Kv.Key); err != nil {
					log.Err(err).Caller().Msg("删除本地数据时出错")
				}
			}
		}
	}
	return nil
}

// 监听存储器数据更新，同步本地数据
func (self *Etcd) watchLoadData(key, value []byte) error {
	keyStr := global.BytesToStr(key)
	// 加载主机
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/hosts/") {
		return self.LoadHost(keyStr, value)
	}
	// 加载上游
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/upstreams/") {
		return self.LoadUpstream(value)
	}
	// 加载路由
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/routes/") {
		return self.LoadRoute(keyStr, value)
	}
	return nil
}

// 监听存储器数据删除，同步本地数据
func (self *Etcd) watchDeleteData(key []byte) error {
	keyStr := global.BytesToStr(key)
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/hosts/") {
		return self.DeleteLocalHost(keyStr)
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/upstreams/") {
		return self.DeleteLocalUpstream(keyStr)
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/routes/") {
		return self.DeleteLocalRoute(keyStr)
	}
	return nil
}
