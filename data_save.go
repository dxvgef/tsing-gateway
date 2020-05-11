package main

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/dxvgef/tsing-gateway/global"
	"go.etcd.io/etcd/clientv3"

	"github.com/rs/zerolog/log"
)

func (p *Proxy) SaveDataToEtcd() (err error) {
	if err = p.SaveUpstreamsToEtcd(); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	if err = p.SaveRoutesToEtcd(); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	if err = p.SaveHostsToEtcd(); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	return nil
}

func (p *Proxy) SaveUpstreamsToEtcd() error {
	var (
		j   []byte
		err error
		key strings.Builder
	)
	// 清空原来的配置
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/upstreams/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = global.EtcdCli.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	key.Reset()

	// 写入upstreams
	for k := range p.Upstreams {
		if j, err = json.Marshal(p.Upstreams[k]); err != nil {
			log.Error().Msg(err.Error())
			return err
		}
		key.WriteString(global.Config.Etcd.KeyPrefix)
		key.WriteString("/upstreams/")
		key.WriteString(p.Upstreams[k].ID)
		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if _, err = global.EtcdCli.Put(ctx2, key.String(), global.BytesToStr(j)); err != nil {
			log.Error().Caller().Msg(err.Error())
			ctx2Cancel()
			return err
		}
		key.Reset()
		ctx2Cancel()
	}
	return nil
}

func (p *Proxy) SaveRoutesToEtcd() error {
	var key strings.Builder

	// 清空原来的配置
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/routes/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err := global.EtcdCli.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	key.Reset()

	// 写入路由
	for routeGroupID, v := range p.Routes {
		for routePath, vv := range v {
			value, err := json.Marshal(vv)
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				return err
			}
			key.WriteString(global.Config.Etcd.KeyPrefix)
			key.WriteString("/routes/")
			key.WriteString(routeGroupID)
			key.WriteString(routePath)
			ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err = global.EtcdCli.Put(ctx2, key.String(), global.BytesToStr(value))
			if err != nil {
				log.Error().Caller().Msg(err.Error())
				ctx2Cancel()
				return err
			}
			key.Reset()
			ctx2Cancel()
		}
	}

	return nil
}

func (p *Proxy) SaveHostsToEtcd() (err error) {
	var key strings.Builder

	// 清空原来的配置
	key.WriteString(global.Config.Etcd.KeyPrefix)
	key.WriteString("/hosts/")
	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer ctxCancel()
	if _, err = global.EtcdCli.Delete(ctx, key.String(), clientv3.WithPrefix()); err != nil {
		log.Error().Caller().Msg(err.Error())
		return err
	}
	key.Reset()

	// 写入路由
	for hostname, upstreamID := range p.Hosts {
		key.WriteString(global.Config.Etcd.KeyPrefix)
		key.WriteString("/hosts/")
		key.WriteString(hostname)
		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = global.EtcdCli.Put(ctx2, key.String(), upstreamID)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			ctx2Cancel()
			return
		}
		key.Reset()
		ctx2Cancel()
	}

	return
}
