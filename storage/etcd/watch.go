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
		if err = proxy.SetHost(path.Base(keyStr), global.BytesToStr(value)); err != nil {
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
		if err = proxy.SetMiddleware(global.BytesToStr(value)); err != nil {
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
		if err = proxy.DelHost(path.Base(keyStr)); err != nil {
			return err
		}
		return nil
	}
	if strings.HasPrefix(keyStr, self.KeyPrefix+"/upstreams/") {
		if err = proxy.DelUpstream(path.Base(keyStr)); err != nil {
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
