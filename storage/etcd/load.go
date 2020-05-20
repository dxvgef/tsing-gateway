package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadAllMiddleware(); err != nil {
		return
	}
	log.Debug().Caller().Interface("middleware", global.Middleware).Msg("加载了middleware")
	if err = self.LoadAllUpstreams(); err != nil {
		return
	}
	log.Debug().Caller().Interface("upstreams", global.Upstreams).Msg("加载了upstreams")
	if err = self.LoadAllRoutes(); err != nil {
		return
	}
	log.Debug().Caller().Interface("routes", global.Routes).Msg("加载了routes")
	if err = self.LoadAllHosts(); err != nil {
		return
	}
	log.Debug().Caller().Interface("hosts", global.Hosts).Msg("加载了hosts")
	return
}

// 加载所有全局中间件
func (self *Etcd) LoadAllMiddleware() error {
	// todo 这个功能还没有实现
	return nil
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
	// 获取路由
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
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
		err = proxy.SetHost(
			strings.TrimPrefix(global.BytesToStr(resp.Kvs[k].Key), "/hosts/"),
			global.BytesToStr(resp.Kvs[k].Value),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
