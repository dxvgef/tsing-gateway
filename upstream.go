package main

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
)

type Configurator struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// 端点信息
type Endpoint struct {
	Addr      string `json:"addr"`       // 地址
	Weight    int    `json:"weight"`     // 权重
	TTL       int    `json:"ttl"`        // 生命周期(秒)
	UpdatedAt int64  `json:"updated_at"` // 最后更新时间
}

type Upstream struct {
	ID         string         `json:"id"`                   // 上游ID
	Middleware []Configurator `json:"middleware,omitempty"` // 中间件配置
	Explorer   Configurator   `json:"explorer"`             // 节点探索器配置
}

func (p *Proxy) NewUpstream(upstream Upstream, persistent bool) error {
	if upstream.ID == "" {
		return errors.New("must specify upstream ID")
	}
	if _, exist := p.Upstreams[upstream.ID]; exist {
		return errors.New("upstream ID:" + upstream.ID + " already exists")
	}
	p.Upstreams[upstream.ID] = upstream
	return nil
}

// set upstream,create if it doesn't exist
func (p *Proxy) SetUpstream(upstream Upstream, persistent bool) error {
	if upstream.ID == "" {
		upstream.ID = global.GetIDStr()
	}
	if upstream.ID == "" {
		return errors.New("must specify upstream ID")
	}
	p.Upstreams[upstream.ID] = upstream
	return nil
}

func (p *Proxy) MatchUpstream(upstreamID string) (upstream Upstream, exist bool) {
	if upstreamID == "" {
		return
	}
	_, exist = p.Upstreams[upstreamID]
	if !exist {
		return
	}
	return p.Upstreams[upstreamID], true
}
