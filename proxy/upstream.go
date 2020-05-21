package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
)

func SetUpstream(upstream global.UpstreamType) error {
	if upstream.ID == "" {
		upstream.ID = global.SnowflakeNode.Generate().String()
	}
	if upstream.ID == "" {
		return errors.New("没有传入upstream.ID,并且无法自动创建ID")
	}
	global.Upstreams[upstream.ID] = upstream
	return nil
}

func DelUpstream(upstreamID string) error {
	if upstreamID == "" {
		return errors.New("upstreamID不能为空")
	}
	delete(global.Upstreams, upstreamID)
	return nil
}

func matchUpstream(upstreamID string) (upstream global.UpstreamType, exist bool) {
	if upstreamID == "" {
		return
	}
	_, exist = global.Upstreams[upstreamID]
	if !exist {
		return
	}
	return global.Upstreams[upstreamID], true
}
