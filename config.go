package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.etcd.io/etcd/clientv3"

	"github.com/dxvgef/tsing-gateway/global"
)

func (p *Proxy) LoadConfigFromJSON(configStr string) error {
	return json.Unmarshal(global.StrToBytes(configStr), p)
}

func (p *Proxy) LoadConfigFromEtcd() error {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	// 获取upstreams
	resp, err := global.EtcdCli.Get(ctx, global.LocalConfig.Etcd.KeyPrefix+"/upstreams/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return err
	}
	for k := range resp.Kvs {
		var upstream Upstream
		err = json.Unmarshal(resp.Kvs[k].Value, &upstream)
		if err != nil {
			log.Debug().Caller().Err(err).Send()
			return err
		}
		err = p.SetUpstream(upstream, false)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
	}
	// 获取路由组
	ctx, ctxCancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	resp, err = global.EtcdCli.Get(ctx, global.LocalConfig.Etcd.KeyPrefix+"/route_groups/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal().Caller().Msg(err.Error())
		return err
	}
	methods := make(map[string]string)
	for k := range resp.Kvs {
		routeGroupID, routePath := parseRouteGroup(resp.Kvs[k].Key)
		err = json.Unmarshal(resp.Kvs[k].Value, &methods)
		if err != nil {
			log.Fatal().Caller().Msg(err.Error())
			return err
		}
		for routeMethod, upstreamID := range methods {
			err = p.SetRoute(routeGroupID, routePath, routeMethod, upstreamID, false)
			if err != nil {
				log.Fatal().Caller().Msg(err.Error())
				return err
			}
		}
	}
	return err
	// return json.Unmarshal(global.StrToBytes(configStr), p)
}

// 从etcd key里解析路由组信息
func parseRouteGroup(key []byte) (routeGroupID, routePath string) {
	keyStr := global.TrimPrefix(key, "/route_groups/")
	pos := strings.Index(keyStr, "/")
	if pos == -1 {
		return
	}
	routeGroupID = keyStr[:pos]
	routePath = strings.TrimLeft(keyStr, routeGroupID)
	return
}
