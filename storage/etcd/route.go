package etcd

import (
	"context"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"

	"github.com/coreos/etcd/clientv3"
)

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

// 存储所有route数据
func (self *Etcd) SaveAllRoutes() (err error) {
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
		err = self.PutRoute(routeGroupID, routePath, routeMethod, v)
		if err != nil {
			return
		}
	}

	return
}

// 设置单个route，如果不存在则创建
func (self *Etcd) PutRoute(routeGroupID, routePath, routeMethod, upstreamID string) error {
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

// 删除单个route
func (self *Etcd) DelRoute(routeGroupID, routePath, routeMethod string) error {
	if routeGroupID == "" {
		return errors.New("路由组ID不能为空")
	}
	if routeMethod != "" {
		routeMethod = strings.ToUpper(routeMethod)
		if !global.InStr(global.HTTPMethods, routeMethod) {
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
