package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
)

func SetUpstream(upstream global.UpstreamType) error {
	if upstream.ID == "" {
		return errors.New("upstream ID不能为空")
	}
	global.Upstreams.Store(upstream.ID, upstream)
	if err := SetUpstreamMiddleware(upstream.ID, upstream.Middleware); err != nil {
		return err
	}
	return nil
}

func DelUpstream(upstreamID string) error {
	if upstreamID == "" {
		return errors.New("upstream ID不能为空")
	}
	global.Upstreams.Delete(upstreamID)
	global.UpstreamMiddleware.Delete(upstreamID)
	return nil
}

func matchUpstream(upstreamID string) (global.UpstreamType, bool) {
	if upstreamID == "" {
		return global.UpstreamType{}, false
	}
	upstream, exist := global.Upstreams.Load(upstreamID)
	if !exist {
		return global.UpstreamType{}, false
	}
	return upstream.(global.UpstreamType), true
}
