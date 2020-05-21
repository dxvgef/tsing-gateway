package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
)

func SetUpstream(upstream global.UpstreamType) error {
	if upstream.ID == "" {
		return errors.New("upstream ID不能为空")
	}
	global.Upstreams[upstream.ID] = upstream
	return nil
}

func DelUpstream(upstreamID string) error {
	if upstreamID == "" {
		return errors.New("upstream ID不能为空")
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
