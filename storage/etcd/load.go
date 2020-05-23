package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadMiddleware(); err != nil {
		return
	}
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

// 加载所有全局中间件
func (self *Etcd) LoadMiddleware() error {
	var str strings.Builder

	// str做为key使用
	str.WriteString(self.KeyPrefix)
	str.WriteString("/middleware")

	// 获取middleware
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, str.String())
	if err != nil {
		return err
	}

	// 重置str准备做为value使用
	str.Reset()
	if resp.Count > 0 {
		str.WriteString(global.BytesToStr(resp.Kvs[0].Value))
	}
	err = proxy.SetGlobalMiddleware(str.String())
	if err != nil {
		return err
	}
	return nil
}

// 加载所有upstream
func (self *Etcd) LoadAllUpstreams() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")

	// 获取upstreams
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		var upstream global.UpstreamType
		err = upstream.UnmarshalJSON(resp.Kvs[k].Value)
		if err != nil {
			return err
		}
		if upstream.ID == "" {
			continue
		}
		err = proxy.SetUpstream(upstream)
		if err != nil {
			return err
		}
	}
	return nil
}

// 加载所有route
func (self *Etcd) LoadAllRoutes() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")

	// 获取路由
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		routeGroupID, routePath, routeMethod, err := global.ParseRoute(global.BytesToStr(resp.Kvs[k].Key), self.KeyPrefix)
		if err != nil {
			return err
		}
		if routeMethod == "" {
			return nil
		}
		err = proxy.SetRoute(routeGroupID, routePath, routeMethod, global.BytesToStr(resp.Kvs[k].Value))
		if err != nil {
			return err
		}
	}
	return nil
}

// 加载所有host
func (self *Etcd) LoadAllHosts() error {
	var hostname string
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		hostname = strings.TrimPrefix(global.BytesToStr(resp.Kvs[k].Key), key.String())
		hostname, err = global.DecodeKey(hostname)
		if err != nil {
			return err
		}
		err = proxy.SetHost(hostname, global.BytesToStr(resp.Kvs[k].Value))
		if err != nil {
			return err
		}
	}
	return nil
}
