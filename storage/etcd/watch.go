package etcd

import (
	"context"
	"encoding/json"
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
				log.Debug().
					Str("type", event.Type.String()).
					Str("key", string(event.Kv.Key)).
					Str("value", string(event.Kv.Value)).
					Send()
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
func (self *Etcd) setRoute(key, value []byte) (err error) {
	routeID, routePath := parseRouteGroup(key)
	log.Debug().Str("routeID", routeID).Str("routePath", routePath).Send()
	methods := make(map[string]string)
	if err = json.Unmarshal(value, &methods); err != nil {
		return
	}
	log.Debug().Interface("methods", methods).Send()
	for method, upstreamID := range methods {
		if err = self.e.SetRoute(routeID, routePath, method, upstreamID); err != nil {
			return
		}
	}
	log.Debug().Caller().Interface("proxy已更新", self.e).Send()
	return
}
