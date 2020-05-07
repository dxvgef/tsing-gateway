package main

import (
	"errors"

	"tsing-gateway/plugin"
)

// 上游信息
type Upstream struct {
	ID        string     `json:"id"`        // 上游ID
	Endpoints []Endpoint `json:"endpoints"` // 端点列表
	Plugins   struct {
		Favicon plugin.Favicon `json:"favicon"` // /favicon请求处理
	} `json:"plugins"`
	// 健康检查
	HealthCheck struct {
		Active struct {
			On         bool   `json:"on"`       // 打开健康检查
			Interval   int    `json:"interval"` // 检查间隔的时间(秒)
			URL        string `json:"url"`      // 主动检查地址
			StatusCode []int  `json:"status_code"`
		} `json:"active"` // 主动检查配置
		Passive struct {
			On  bool `json:"on"`  // 打开健康检查
			TTL int  `json:"ttl"` // 端点的生命周期(秒)
		} `json:"passive"` // 被动检查配置
	} `json:"health_check"`
}

// 端点信息
type Endpoint struct {
	Addr      string `json:"addr"`       // 地址
	Weight    int    `json:"weight"`     // 权重
	TTL       int    `json:"ttl"`        // 生命周期(秒)
	UpdatedAt int64  `json:"updated_at"` // 最后更新时间
}

// 新建上游
func (p *Proxy) newUpstream(upstream Upstream, persistent bool) error {
	if upstream.ID == "" {
		return errors.New("没有传入上游ID")
	}
	if _, exist := p.upstreams[upstream.ID]; exist {
		return errors.New("上游ID:" + upstream.ID + "已存在")
	}
	p.upstreams[upstream.ID] = upstream
	return nil
}

// 设置上游，如果存在则更新，不存在则创建
func (p *Proxy) setUpstream(upstream Upstream, persistent bool) error {
	if upstream.ID == "" {
		upstream.ID = getID()
	}
	if upstream.ID == "" {
		return errors.New("上游ID不能为空")
	}
	p.upstreams[upstream.ID] = upstream
	return nil
}

// 匹配上游
func (p *Proxy) matchUpstream(upstreamID string) (upstream Upstream, exist bool) {
	if upstreamID == "" {
		return
	}
	_, exist = p.upstreams[upstreamID]
	if !exist {
		return
	}
	return p.upstreams[upstreamID], true
}
