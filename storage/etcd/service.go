package etcd

import (
	"context"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"

	"github.com/coreos/etcd/clientv3"
)

// 从存储器加载服务数据到本地
func (self *Etcd) LoadService(data []byte) error {
	var service global.ServiceType
	err := service.UnmarshalJSON(data)
	if err != nil {
		log.Err(err).Caller().Msg("加载服务时出错")
		return err
	}
	return proxy.SetService(service)
}

// 从存储器加载所有服务数据到本地
func (self *Etcd) LoadAllService() error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")

	// 获取services
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err := self.client.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for k := range resp.Kvs {
		err = self.LoadService(resp.Kvs[k].Value)
		if err != nil {
			log.Err(err).Caller().Msg("加载所有服务时出错")
			return err
		}
	}
	return nil
}

// 将本地服务数据保存到存储器
func (self *Etcd) SaveService(serviceID, config string) error {
	serviceID = global.EncodeKey(serviceID)
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	key.WriteString(serviceID)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), config); err != nil {
		return err
	}
	return nil
}

// 将本地所有服务数据保存到存储器
func (self *Etcd) SaveAllService() error {
	var (
		err         error
		key         strings.Builder
		services    = make(map[string]string, global.SyncMapLen(&global.Services))
		configBytes []byte
	)

	// 将数据保存到临时变量中
	global.Services.Range(func(k, v interface{}) bool {
		service, ok := v.(global.ServiceType)
		if !ok {
			log.Err(err).Caller().Msg("服务配置的类型断言失败")
			return false
		}
		if configBytes, err = json.Marshal(&service); err != nil {
			log.Err(err).Caller().Msg("服务配置序列化成JSON失败")
			return false
		}
		services[k.(string)] = global.BytesToStr(configBytes)
		return true
	})

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Msg("清空存储器中的服务数据失败")
		return err
	}

	// 将内存中的数据写入到存储器中
	for k := range services {
		if err = self.SaveService(k, services[k]); err != nil {
			return err
		}
	}
	return nil
}

// 删除本地服务数据
func (self *Etcd) DeleteLocalService(key string) error {
	serviceID, err := global.DecodeKey(path.Base(key))
	if err != nil {
		return err
	}
	return proxy.DelService(serviceID)
}

// 删除本地服务数据
func (self *Etcd) DeleteStorageService(serviceID string) error {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/services/")
	key.WriteString(serviceID)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String())
	if err != nil {
		log.Err(err).Caller().Msg("删除存储器中的服务数据失败")
		return err
	}
	return nil
}
