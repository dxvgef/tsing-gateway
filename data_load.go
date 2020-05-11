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

func (p *Proxy) DataToJSON() (string, error) {
	j, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return global.BytesToStr(j), nil
}

func (p *Proxy) LoadDataFromJSON(configStr string) error {
	return json.Unmarshal(global.StrToBytes(configStr), p)
}

func (p *Proxy) LoadDataFromEtcd() error {
	var err error
	if err = p.LoadUpstreamsFromEtcd(); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}

	log.Debug().Interface("加载配置", p).Send()
	return nil
}

// 从etcd key里解析路由组信息
func parseRouteGroup(key []byte) (routeGroupID, routePath string) {
	keyStr := global.TrimPrefix(key, "/routes/")
	pos := strings.Index(keyStr, "/")
	if pos == -1 {
		return
	}
	routeGroupID = keyStr[:pos]
	routePath = strings.TrimLeft(keyStr, routeGroupID)
	return
}

func (p *Proxy) LoadUpstreamsFromEtcd() error {
	var key strings.Builder
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	// 获取upstreams
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/upstreams/")
	resp, err := global.EtcdCli.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	for k := range resp.Kvs {
		var upstream Upstream
		err = json.Unmarshal(resp.Kvs[k].Value, &upstream)
		if err != nil {
			return err
		}
		err = p.SetUpstream(upstream, false)
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
	}
	return nil
}

func (p *Proxy) LoadRoutesFromEtcd() error {
	var key strings.Builder
	// 获取路由
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/routes/")
	resp, err := global.EtcdCli.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	methods := make(map[string]string)
	for k := range resp.Kvs {
		routeGroupID, routePath := parseRouteGroup(resp.Kvs[k].Key)
		err = json.Unmarshal(resp.Kvs[k].Value, &methods)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return err
		}
		for routeMethod, upstreamID := range methods {
			err = p.SetRoute(routeGroupID, routePath, routeMethod, upstreamID, false)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return err
			}
		}
	}
	return nil
}

// 获取主机
func (p *Proxy) LoadHostsFromEtcd() error {
	var key strings.Builder
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/hosts/")
	resp, err := global.EtcdCli.Get(ctx, key.String(), clientv3.WithPrefix())
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	for k := range resp.Kvs {
		err = p.SetHost(
			global.TrimPrefix(resp.Kvs[k].Key, "/hosts/"),
			global.BytesToStr(resp.Kvs[k].Value),
			false,
		)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return err
		}
	}
	return nil
}
