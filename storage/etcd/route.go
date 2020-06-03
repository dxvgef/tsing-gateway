package etcd

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"

	"github.com/coreos/etcd/clientv3"
)

// 从存储器加载路由到本地
func (self *Etcd) LoadRoute(key string, data []byte) error {
	routeGroupID, routePath, routeMethod, err := global.ParseRoute(key, self.KeyPrefix)
	if err != nil {
		return err
	}
	if !global.InStr(global.HTTPMethods, routeMethod) {
		return errors.New("HTTP方法无效")
	}

	return proxy.SetRoute(routeGroupID, routePath, routeMethod, global.BytesToStr(data))
}

// 从存储器加载所有路由到本地
func (self *Etcd) LoadAllRoute() error {
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

	// 将所有路由都加载到缓存
	for k := range resp.Kvs {
		if err = self.LoadRoute(global.BytesToStr(resp.Kvs[k].Key), resp.Kvs[k].Value); err != nil {
			return err
		}
	}
	return nil
}

// 保存本地路由到存储器，如果不存在则创建
func (self *Etcd) SaveRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
	routeMethod = strings.ToUpper(routeMethod)
	if !global.InStr(global.HTTPMethods, routeMethod) {
		return errors.New("HTTP方法无效")
	}

	routeGroupID = global.EncodeKey(routeGroupID)
	routePath = global.EncodeKey(routePath)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
	key.WriteString(routeGroupID)
	key.WriteString("@")
	key.WriteString(routePath)
	key.WriteString("@")
	key.WriteString(routeMethod)
	path.Join()

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), upstreamID); err != nil {
		return err
	}
	return nil
}

// 将本地所有路由保存到存储器
func (self *Etcd) SaveAllRoute() (err error) {
	var (
		key    strings.Builder
		routes = make(map[string]string, global.SyncMapLen(&global.Routes))
	)

	// 将数据保存到临时变量中
	global.Routes.Range(func(k, v interface{}) bool {
		routes[k.(string)] = v.(string)
		return true
	})

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/routes/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		return err
	}

	// 将内存中的数据写入到存储器中
	var routeGroupID, routePath, routeMethod string
	for k, v := range routes {
		routeGroupID, routePath, routeMethod, err = global.ParseRoute(k, "")
		if err != nil {
			return
		}
		err = self.SaveRoute(routeGroupID, routePath, routeMethod, v)
		if err != nil {
			return
		}
	}

	return
}

// 删除本地路由数据
func (self *Etcd) DeleteLocalRoute(keyStr string) error {
	routeGroupID, routePath, routeMethod, err := global.ParseRoute(keyStr, self.KeyPrefix)
	if err != nil {
		return err
	}
	return proxy.DeleteRoute(routeGroupID, routePath, routeMethod)
}

// 删除存储器中路由数据
func (self *Etcd) DeleteStorageRoute(routeGroupID, routePath, routeMethod string) error {
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
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String())
	if err != nil {
		log.Err(err).Caller().Msg("删除存储器中的路由数据失败")
		return err
	}
	return nil
}
