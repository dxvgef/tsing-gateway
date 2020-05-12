package etcd

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
)

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	if err = self.SaveAllHosts(); err != nil {
		return
	}
	if err = self.SaveAllUpstreams(); err != nil {
		return
	}
	if err = self.SaveAllRoutes(); err != nil {
		return
	}
	return
}

// 存储所有upstream数据
func (self *Etcd) SaveAllUpstreams() error {
	var (
		jsonBytes []byte
		key       strings.Builder
	)

	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	key.Reset()

	// 写入upstreams
	for k := range self.e.Upstreams {
		if jsonBytes, err = self.e.Upstreams[k].MarshalJSON(); err != nil {
			return err
		}
		key.WriteString(self.KeyPrefix)
		key.WriteString("/upstreams/")
		key.WriteString(self.e.Upstreams[k].ID)
		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err = self.client.Put(ctx2, key.String(), global.BytesToStr(jsonBytes)); err != nil {
			ctx2Cancel()
			return err
		}
		key.Reset()
		ctx2Cancel()
	}
	return nil
}

// 存储所有route数据
func (self *Etcd) SaveAllRoutes() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		return err
	}
	key.Reset()

	// 写入路由
	for routeGroupID, v := range self.e.Routes {
		for routePath, vv := range v {
			value, err := json.Marshal(vv)
			if err != nil {
				return err
			}
			key.WriteString(self.KeyPrefix)
			key.WriteString("/routes/")
			key.WriteString(routeGroupID)
			key.WriteString(routePath)
			ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err = self.client.Put(ctx2, key.String(), global.BytesToStr(value))
			if err != nil {
				ctx2Cancel()
				return err
			}
			key.Reset()
			ctx2Cancel()
		}
	}

	return nil
}

// 存储所有host数据
func (self *Etcd) SaveAllHosts() error {
	var key strings.Builder

	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	key.Reset()

	// 写入路由
	for hostname, upstreamID := range self.e.Hosts {
		key.WriteString(self.KeyPrefix)
		key.WriteString("/hosts/")
		key.WriteString(hostname)
		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = self.client.Put(ctx2, key.String(), upstreamID)
		if err != nil {
			ctx2Cancel()
			return err
		}
		key.Reset()
		ctx2Cancel()
	}

	return nil
}
