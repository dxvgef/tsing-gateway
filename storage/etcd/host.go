package etcd

import (
	"context"
	json "encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"
)

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

// 将内存中所有的host数据写入到存储器
func (self *Etcd) SaveAllHosts() error {
	var (
		err         error
		key         strings.Builder
		hosts       = make(map[string]string, global.SyncMapLen(&global.Hosts))
		configBytes []byte
	)

	// 将配置保存到临时变量中
	global.Hosts.Range(func(k, v interface{}) bool {
		h, ok := v.([]global.HostType)
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
		if err = self.PutHost(hostname, config); err != nil {
			return err
		}
	}

	return nil
}

// 设置单个host，如果不存在则创建
func (self *Etcd) PutHost(hostname, config string) error {
	hostname = global.EncodeKey(hostname)

	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/hosts/")
	key.WriteString(hostname)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), config); err != nil {
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
