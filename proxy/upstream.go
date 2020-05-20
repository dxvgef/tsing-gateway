package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
)

type Upstream struct {
	ID         string                `json:"id"`                   // 上游ID
	Middleware []global.ModuleConfig `json:"middleware,omitempty"` // 中间件配置
	Discover   global.ModuleConfig   `json:"discover"`             // 节点发现配置
	// 启用缓存，如果关闭，则每次请求都从etcd中获取endpoints
	Cache bool `json:"cache"`
	/*
		缓存重试次数
		在缓存中失败达到指定次数后，重新从discover中获取endpoints来更新缓存
	*/
	CacheRetry   int               `json:"cache_retry"`
	Endpoints    []global.Endpoint `json:"-"`                      // 终端列表
	LoadBalance  string            `json:"load_balance,omitempty"` // 负载均衡算法
	LastEndpoint string            `json:"-"`                      // 最后使用的endpoint，用于防止连续命中同一个
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

func (p *Engine) matchUpstream(upstreamID string) (upstream Upstream, exist bool) {
	if upstreamID == "" {
		return
	}
	_, exist = p.Upstreams[upstreamID]
	if !exist {
		return
	}
	return p.Upstreams[upstreamID], true
}
