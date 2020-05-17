package proxy

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
	Discover   Configurator   `json:"discover"`             // 节点发现配置
}

func (p *Engine) NewUpstream(upstream Upstream) error {
	if upstream.ID == "" {
		upstream.ID = global.GetIDStr()
	}
	if upstream.ID == "" {
		return errors.New("没有传入upstream.ID,并且无法自动创建ID")
	}
	if _, exist := p.Upstreams[upstream.ID]; exist {
		return errors.New("upstream.ID:" + upstream.ID + " 已存在")
	}
	p.Upstreams[upstream.ID] = upstream
	return nil
}

func (p *Engine) SetUpstream(upstream Upstream) error {
	if upstream.ID == "" {
		upstream.ID = global.GetIDStr()
	}
	if upstream.ID == "" {
		return errors.New("没有传入upstream.ID,并且无法自动创建ID")
	}
	p.Upstreams[upstream.ID] = upstream
	return nil
}

func (p *Engine) DelUpstream(upstreamID string) error {
	if upstreamID == "" {
		return errors.New("upstreamID不能为空")
	}
	delete(p.Upstreams, upstreamID)
	return nil
}

func (p *Engine) MatchUpstream(upstreamID string) (upstream Upstream, exist bool) {
	if upstreamID == "" {
		return
	}
	_, exist = p.Upstreams[upstreamID]
	if !exist {
		return
	}
	return p.Upstreams[upstreamID], true
}
