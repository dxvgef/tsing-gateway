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

// 从存储器加载上游数据到本地
func (self *Etcd) LoadUpstream(data []byte) error {
	var upstream global.UpstreamType
	err := upstream.UnmarshalJSON(data)
	if err != nil {
		log.Err(err).Caller().Msg("加载上游时出错")
		return err
	}
	return proxy.SetUpstream(upstream)
}

// 从存储器加载所有上游数据到本地
func (self *Etcd) LoadAllUpstream() error {
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
		err = self.LoadUpstream(resp.Kvs[k].Value)
		if err != nil {
			log.Err(err).Caller().Msg("加载所有上游时出错")
			return err
		}
	}
	return nil
}

// 将本地上游数据保存到存储器
func (self *Etcd) SaveUpstream(upstreamID, config string) error {
	upstreamID = global.EncodeKey(upstreamID)
	var key strings.Builder
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")
	key.WriteString(upstreamID)

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := self.client.Put(ctx, key.String(), config); err != nil {
		return err
	}
	return nil
}

// 将本地所有上游数据保存到存储器
func (self *Etcd) SaveAllUpstream() error {
	var (
		err         error
		key         strings.Builder
		upstreams   = make(map[string]string, global.SyncMapLen(&global.Upstreams))
		configBytes []byte
	)

	// 将数据保存到临时变量中
	global.Upstreams.Range(func(k, v interface{}) bool {
		upstream, ok := v.(global.UpstreamType)
		if !ok {
			log.Err(err).Caller().Msg("上游配置的类型断言失败")
			return false
		}
		if configBytes, err = json.Marshal(&upstream); err != nil {
			log.Err(err).Caller().Msg("上游配置序列化成JSON失败")
			return false
		}
		upstreams[k.(string)] = global.BytesToStr(configBytes)
		return true
	})

	// 清空存储器中的配置
	key.WriteString(self.KeyPrefix)
	key.WriteString("/upstreams/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	_, err = self.client.Delete(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Err(err).Caller().Msg("清空存储器中的上游数据失败")
		return err
	}

	// 将内存中的数据写入到存储器中
	for k := range upstreams {
		if err = self.SaveUpstream(k, upstreams[k]); err != nil {
			return err
		}
	}
	return nil
}

// 删除本地上游数据
func (self *Etcd) DeleteLocalUpstream(key string) error {
	upstreamID, err := global.DecodeKey(path.Base(key))
	if err != nil {
		return err
	}
	return proxy.DelUpstream(upstreamID)
}
