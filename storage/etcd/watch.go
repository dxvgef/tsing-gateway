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
				if err := self.putData(event.Kv.Key, event.Kv.Value); err != nil {
					log.Err(err).Caller().Send()
				}
			case clientv3.EventTypeDelete:
				if err := self.delData(event.Kv.Key, event.Kv.Value); err != nil {
					log.Err(err).Caller().Send()
				}
			}
		}
	}
	return nil
}

// put数据的操作
func (self *Etcd) putData(key, value []byte) error {
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
				if err = self.setHost(keyStr, value); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/upstreams/":
				if err = self.setUpstream(value); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/routes/":
				log.Debug().Caller().Msg("put了路由")
				if err = self.setRoute(key, value); err != nil {
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

// 设置单个host
func (self *Etcd) setHost(key string, value []byte) (err error) {
	err = self.e.SetHost(path.Base(key), global.BytesToStr(value))
	if err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 设置单个upstream
func (self *Etcd) setUpstream(value []byte) (err error) {
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

// 设置单个route
func (self *Etcd) setRoute(key, value []byte) error {
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
func (self *Etcd) delData(key, value []byte) error {
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
				if err = self.delHost(keyStr); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/upstreams/":
				if err = self.delUpstream(keyStr); err != nil {
					log.Err(err).Caller().Send()
					return err
				}
			case "/routes/":
				if err = self.delRoute(key, value); err != nil {
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

// 删除单个host
func (self *Etcd) delHost(key string) (err error) {
	err = self.e.DelHost(path.Base(key))
	if err != nil {
		return err
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 删除单个upstream
func (self *Etcd) delUpstream(key string) (err error) {
	if err = self.e.DelUpstream(path.Base(key)); err != nil {
		return
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}

// 删除单个route
func (self *Etcd) delRoute(key, value []byte) error {
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
