package etcd

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadAllUpstreams(); err != nil {
		return
	}
	if err = self.LoadAllRoutes(); err != nil {
		return
	}
	if err = self.LoadAllHosts(); err != nil {
		return
	}
	return
}

// 加载所有upstream
func (self *Etcd) LoadAllUpstreams() error {
	var key strings.Builder
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	// 获取upstreams
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		var upstream proxy.Upstream
		err = upstream.UnmarshalJSON(resp.Kvs[k].Value)
		if err != nil {
			return err
		}
		err = self.e.SetUpstream(upstream)
		if err != nil {
			return err
		}
	}
	return nil
}

// 加载所有route
func (self *Etcd) LoadAllRoutes() error {
	var key strings.Builder
	// 获取路由
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	methods := make(map[string]string)
	for k := range resp.Kvs {
		routeGroupID, routePath := parseRouteGroup(resp.Kvs[k].Key)
		err = json.Unmarshal(resp.Kvs[k].Value, &methods)
		if err != nil {
			return err
		}
		for routeMethod, upstreamID := range methods {
			err = self.e.SetRoute(routeGroupID, routePath, routeMethod, upstreamID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 加载所有host
func (self *Etcd) LoadAllHosts() error {
	var key strings.Builder
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		err = self.e.SetHost(
			global.TrimPrefix(resp.Kvs[k].Key, "/hosts/"),
			global.BytesToStr(resp.Kvs[k].Value),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// 从etcd key里解析路由组信息
func parseRouteGroup(key []byte) (routeGroupID, routePath string) {
	keyStr := global.TrimPrefix(key, "/routes/")
	pos := strings.Index(keyStr, "/")
	if pos == -1 {
		return
	}
	routeGroupID = keyStr[:pos]
	routePath = strings.TrimLeft(keyStr, routeGroupID)
	return
}
