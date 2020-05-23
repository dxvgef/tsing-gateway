package etcd

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
)

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	if err = self.SaveMiddleware(); err != nil {
		return
	}
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

// 存储全局中间件数据
func (self *Etcd) SaveMiddleware() error {
	return nil
}

// 设置全局中间件数据
func (self *Etcd) PutMiddleware(configStr string) error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/middleware")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), configStr); err != nil {
		return err
	}
	return nil
}

// 存储所有upstream数据
func (self *Etcd) SaveAllUpstreams() error {
	var (
		err       error
		jsonBytes []byte
		key       strings.Builder
	)

	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")

	// 将配置保存到临时变量中
	upstreams := make(map[string]string, len(global.Upstreams))
	for k := range global.Upstreams {
		if k == "" {
			continue
		}
		jsonBytes, err = global.Upstreams[k].MarshalJSON()
		if err != nil {
			continue
		}
		upstreams[k] = global.BytesToStr(jsonBytes)
	}

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 写入upstreams
	for k := range upstreams {
		key.Reset()
		key.WriteString(self.KeyPrefix)
		key.WriteString("/upstreams/")
		key.WriteString(global.EncodeKey(k))
		ctxTmp, ctxTmpCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err = self.client.Put(ctxTmp, key.String(), upstreams[k]); err != nil {
			ctxTmpCancel()
			return err
		}
		ctxTmpCancel()
	}
	return nil
}

// 存储所有route数据
func (self *Etcd) SaveAllRoutes() (err error) {
	var (
		key    strings.Builder
		routes = make(map[string]map[string]map[string]string)
	)

	// 将配置保存到临时变量中
	for routeGroupID, v := range global.Routes {
		if _, exist := routes[routeGroupID]; !exist {
			routes[routeGroupID] = make(map[string]map[string]string)
		}
		for routePath, vv := range v {
			if _, exist := routes[routeGroupID][routePath]; !exist {
				routes[routeGroupID][routePath] = make(map[string]string)
			}
			for routeMethod, upstreamID := range vv {
				if routeMethod == "" {
					continue
				}
				routes[routeGroupID][routePath][routeMethod] = upstreamID
			}
		}
	}

	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		return err
	}

	// 写入路由
	for routeGroupID, v := range routes {
		for routePath, vv := range v {
			for routeMethod, upstreamID := range vv {
				if routeMethod == "" {
					continue
				}
				key.Reset()
				key.WriteString(self.KeyPrefix)
				key.WriteString("/routes/")
				key.WriteString(global.EncodeKey(routeGroupID))
				key.WriteString("/")
				key.WriteString(global.EncodeKey(routePath))
				key.WriteString("/")
				key.WriteString(routeMethod)

				ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_, err = self.client.Put(ctx2, key.String(), upstreamID)
				if err != nil {
					ctx2Cancel()
					return
				}
				ctx2Cancel()
			}
		}
	}

	return
}

// 存储所有host数据
func (self *Etcd) SaveAllHosts() error {
	var (
		key   strings.Builder
		hosts = make(map[string]string)
	)

	// 将配置保存到临时变量中
	for hostname, upstreamID := range global.Hosts {
		hosts[hostname] = upstreamID
	}

	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")

	// 清空原来的配置
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}

	// 写入路由
	for hostname, upstreamID := range hosts {
		key.Reset()
		key.WriteString(self.KeyPrefix)
		key.WriteString("/hosts/")
		key.WriteString(global.EncodeKey(hostname))
		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = self.client.Put(ctx2, key.String(), upstreamID)
		if err != nil {
			ctx2Cancel()
			return err
		}
		ctx2Cancel()
	}

	return nil
}

// 设置单个host，如果不存在则创建
func (self *Etcd) PutHost(hostname, upstreamID string) error {
	hostname = global.EncodeKey(hostname)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	key.WriteString(hostname)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), upstreamID); err != nil {
		return err
	}
	return nil
}

// 删除host
func (self *Etcd) DelHost(hostname string) error {
	hostname = global.EncodeKey(hostname)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	key.WriteString(hostname)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Delete(ctx, key.String()); err != nil {
		return err
	}
	return nil
}

// 设置单个upstream，如果不存在则创建
func (self *Etcd) PutUpstream(upstreamID, upstreamConfig string) error {
	upstreamID = global.EncodeKey(upstreamID)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")
	key.WriteString(upstreamID)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), upstreamConfig); err != nil {
		return err
	}
	return nil
}

// 删除upstream
func (self *Etcd) DelUpstream(upstreamID string) error {
	upstreamID = global.EncodeKey(upstreamID)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")
	key.WriteString(upstreamID)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Delete(ctx, key.String()); err != nil {
		return err
	}
	return nil
}

// 设置单个route，如果不存在则创建
func (self *Etcd) PutRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
	routeMethod = strings.ToUpper(routeMethod)
	if !global.InStr(global.Methods, routeMethod) {
		return errors.New("HTTP方法无效")
	}

	routeGroupID = global.EncodeKey(routeGroupID)
	routePath = global.EncodeKey(routePath)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
	key.WriteString(routeGroupID)
	key.WriteString("/")
	key.WriteString(routePath)
	key.WriteString("/")
	key.WriteString(routeMethod)
	path.Join()

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), upstreamID); err != nil {
		return err
	}
	return nil
}

// 删除单个route
func (self *Etcd) DelRoute(routeGroupID, routePath, routeMethod string) error {
	if routeGroupID == "" {
		return errors.New("路由组ID不能为空")
	}
	if routeMethod != "" {
		routeMethod = strings.ToUpper(routeMethod)
		if !global.InStr(global.Methods, routeMethod) {
			return errors.New("HTTP方法无效")
		}
	}

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
	routeGroupID = global.EncodeKey(routeGroupID)
	key.WriteString(routeGroupID)
	key.WriteString("/")
	routePath = global.EncodeKey(routePath)
	key.WriteString(routePath)
	key.WriteString("/")

	if routeMethod != "" {
		key.WriteString(routeMethod)
	}

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if routeMethod == "" {
		if _, err := self.client.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
			return err
		}
		return nil
	}
	if _, err := self.client.Delete(ctx, key.String()); err != nil {
		return err
	}
	return nil
}
