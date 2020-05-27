package etcd

import (
	"context"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
)

// 获取所有Endpoint列表
func (self *Etcd) Fetch(upstreamID string) ([]global.EndpointType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	resp, err := self.client.Get(ctx, upstreamID+"/", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, nil
	}
	eps := make([]global.EndpointType, resp.Count)
	for k := range resp.Kvs {
		var ep global.EndpointType
		ep.UpstreamID = upstreamID
		ep.Addr, err = global.DecodeKey(global.BytesToStr(resp.Kvs[k].Key))
		if err != nil {
			return nil, err
		}
		ep.Weight, err = strconv.Atoi(global.BytesToStr(resp.Kvs[k].Value))
		if err != nil {
			return nil, err
		}
		eps[k] = ep
	}
	return eps, nil
}
