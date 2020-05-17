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
					log.Err(err).Caller().Send()
				}
			case clientv3.EventTypeDelete:
				if err := self.delDataToLocal(event.Kv.Key); err != nil {
					log.Err(err).Caller().Send()
				}
			}
		}
	}
	return nil
}

// put数据的操作
func (self *Etcd) putDataToLocal(key, value []byte) error {
	var (
		err     error
		keyStr  = global.BytesToStr(key)
		modules = []string{"/hosts/", "/upstreams/", "/routes/"}
		keyPath strings.Builder
	)
	for k := range modules {
		keyPath.WriteString(self.KeyPrefix)
		keyPath.WriteString(modules[k])
		if strings.HasPrefix(keyStr, keyPath.String()) {
			switch modules[k] {
			case "/hosts/":
				if err = self.setHostToLocal(keyStr, value); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/upstreams/":
				if err = self.setUpstreamToLocal(value); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/routes/":
				log.Debug().Caller().Msg("put了路由")
				if err = self.setRouteToLocal(key, value); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			}
			break
		}
		keyPath.Reset()
	}
	return nil
}

// 设置本地单个host
func (self *Etcd) setHostToLocal(key string, value []byte) (err error) {
	err = self.e.SetHost(path.Base(key), global.BytesToStr(value))
	if err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 设置本地单个upstream
func (self *Etcd) setUpstreamToLocal(value []byte) (err error) {
	var upstream proxy.Upstream
	if err = upstream.UnmarshalJSON(value); err != nil {
		return
	}
	if err = self.e.SetUpstream(upstream); err != nil {
		return
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 设置本地单个route
func (self *Etcd) setRouteToLocal(key, value []byte) error {
	routeID, routePath, routeMethod, err := parseRouteGroup(key)
	if err != nil {
		return err
	}
	if routeMethod == "" {
		return nil
	}
	if err = self.e.SetRoute(routeID, routePath, routeMethod, global.BytesToStr(value)); err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return nil
}

// del数据的操作
func (self *Etcd) delDataToLocal(key []byte) error {
	var (
		err     error
		keyStr  = global.BytesToStr(key)
		modules = []string{"/hosts/", "/upstreams/", "/routes/"}
		keyPath strings.Builder
	)
	for k := range modules {
		keyPath.WriteString(self.KeyPrefix)
		keyPath.WriteString(modules[k])
		if strings.HasPrefix(keyStr, keyPath.String()) {
			switch modules[k] {
			case "/hosts/":
				if err = self.delHostToLocal(keyStr); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/upstreams/":
				if err = self.delUpstreamToLocal(keyStr); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/routes/":
				if err = self.delRouteToLocal(key); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			}
			break
		}
		keyPath.Reset()
	}
	return nil
}

// 删除本地单个host
func (self *Etcd) delHostToLocal(key string) (err error) {
	err = self.e.DelHost(path.Base(key))
	if err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 删除本地单个upstream
func (self *Etcd) delUpstreamToLocal(key string) (err error) {
	if err = self.e.DelUpstream(path.Base(key)); err != nil {
		return
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 删除本地单个route
func (self *Etcd) delRouteToLocal(key []byte) error {
	log.Debug().Caller().Msg(global.BytesToStr(key))
	routeID, routePath, routeMethod, err := parseRouteGroup(key)
	if err != nil {
		return err
	}
	if err = self.e.DelRoute(routeID, routePath, routeMethod); err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return nil
}
