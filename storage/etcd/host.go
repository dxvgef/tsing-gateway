package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

// 从存储器加载主机数据到本地，如果不存在则创建
func (self *Etcd) LoadHost(key string, data []byte) error {
	hostname, err := global.DecodeKey(path.Base(key))
	if err != nil {
		return err
	}
	var host global.HostType
	err = host.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	return proxy.SetHost(hostname, host)
}

// 从存储器加载所有主机数据到本地
func (self *Etcd) LoadAllHost() error {
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
		err = self.LoadHost(global.BytesToStr(resp.Kvs[k].Key), resp.Kvs[k].Value)
		if err != nil {
			return err
		}
	}
	return nil
}

// 将本地主机数据保存到存储器中，如果不存在则创建
func (self *Etcd) SaveHost(hostname, config string) (err error) {
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	key.WriteString(hostname)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = self.client.Put(ctx, key.String(), config); err != nil {
		return err
	}
	return nil
}

// 将本地中所有主机数据保存到存储器
func (self *Etcd) SaveAllHost() error {
	var (
		err         error
		key         strings.Builder
		hosts       = make(map[string]string, global.SyncMapLen(&global.Hosts))
		configBytes []byte
	)

	// 将配置保存到临时变量中
	global.Hosts.Range(func(k, v interface{}) bool {
		h, ok := v.(global.HostType)
		if !ok {
			err = errors.New("主机" + k.(string) + "的配置异常")
			log.Err(err).Caller().Msg("主机配置的类型断言失败")
			return false
		}
		if configBytes, err = json.Marshal(&h); err != nil {
			log.Err(err).Caller().Msg("主机配置序列化成JSON失败")
			return false
		}
		hosts[k.(string)] = global.BytesToStr(configBytes)
		return true
	})
	if err != nil {
		return err
	}

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Msg("清空存储器中的上游数据失败")
		return err
	}

	// 将内存中的数据写入到存储器中
	for hostname, config := range hosts {
		hostname = global.EncodeKey(hostname)
		if err = self.SaveHost(hostname, config); err != nil {
			return err
		}
	}

	return nil
}

// 删除本地主机数据
func (self *Etcd) DeleteLocalHost(key string) error {
	hostname, err := global.DecodeKey(path.Base(key))
	if err != nil {
		return err
	}
	return proxy.DelHost(hostname)
}

// 删除存储器中主机数据
func (self *Etcd) DeleteStorageHost(hostname string) error {
	if hostname == "" {
		return errors.New("主机名不能为空")
	}

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	key.WriteString(hostname)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err := self.client.Delete(ctx, key.String())
	if err != nil {
		log.Err(err).Caller().Msg("删除存储器中的主机数据失败")
		return err
	}
	return nil
}
