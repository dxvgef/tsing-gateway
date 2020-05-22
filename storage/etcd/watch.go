package etcd

import (
	"context"
	"path"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

// 监听变更
func (self *Etcd) Watch() error {
	ch := self.client.Watch(context.Background(), self.KeyPrefix+"/", clientv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			switch event.Type {
			// 添加key
			case clientv3.EventTypePut:
				if err := self.putDataToLocal(event.Kv.Key, event.Kv.Value); err != nil {
					log.Err(err).Caller().Msg("更新本地数据时出错")
				}
			case clientv3.EventTypeDelete:
				if err := self.delDataToLocal(event.Kv.Key); err != nil {
					log.Err(err).Caller().Msg("删除本地数据时出错")
				}
			}
		}
	}
	return nil
}

// put数据的操作
func (self *Etcd) putDataToLocal(key, value []byte) error {
	var (
		err    error
		keyStr = global.BytesToStr(key)
	)
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/hosts/") {
		var hostname string
		hostname, err = global.DecodeKey(path.Base(keyStr))
		if err != nil {
			return err
		}
		if err = proxy.SetHost(hostname, global.BytesToStr(value)); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/upstreams/") {
		var upstream global.UpstreamType
		if err = upstream.UnmarshalJSON(value); err != nil {
			return err
		}
		if err = proxy.SetUpstream(upstream); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/routes/") {
		// todo 这里的parseRoute函数好像没有正确的解码key，所以数据都写不进去
		routeID, routePath, routeMethod, errTmp := global.ParseRoute(keyStr, self.KeyPrefix)
		if errTmp != nil {
			return err
		}
		if err = proxy.SetRoute(routeID, routePath, routeMethod, global.BytesToStr(value)); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/middleware") {
		if err = proxy.SetGlobalMiddleware(global.BytesToStr(value)); err != nil {
			log.Err(err).Caller().Msg("更新本地路由数据时出错")
			return err
		}
		return nil
	}
	return nil
}

// del数据的操作
func (self *Etcd) delDataToLocal(key []byte) error {
	var (
		err    error
		keyStr = global.BytesToStr(key)
	)
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/hosts/") {
		var mapKey string
		mapKey, err = global.DecodeKey(path.Base(keyStr))
		if err != nil {
			return err
		}
		if err = proxy.DelHost(mapKey); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/upstreams/") {
		var mapKey string
		mapKey, err = global.DecodeKey(path.Base(keyStr))
		if err != nil {
			return err
		}
		if err = proxy.DelUpstream(mapKey); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/routes/") {
		routeID, routePath, routeMethod, err := global.ParseRoute(keyStr, self.KeyPrefix)
		if err != nil {
			return err
		}
		if err = proxy.DelRoute(routeID, routePath, routeMethod); err != nil {
			return err
		}
		return nil
	}
	return nil
}
