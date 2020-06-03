package etcd

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/proxy"

	"github.com/coreos/etcd/clientv3"
)

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

// 存储所有upstream数据
func (self *Etcd) SaveAllUpstreams() error {
	var (
		err         error
		key         strings.Builder
		upstreams   = make(map[string]string, global.SyncMapLen(&global.Upstreams))
		configBytes []byte
	)

	// 将数据保存到临时变量中
	global.Upstreams.Range(func(k, v interface{}) bool {
		u, ok := v.(global.UpstreamType)
		if !ok {
			log.Err(err).Caller().Msg("上游配置的类型断言失败")
			return false
		}
		configBytes, err = u.MarshalJSON()
		if err != nil {
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
