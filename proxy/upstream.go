package proxy

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/middleware"
)

func SetUpstream(upstream global.UpstreamType) error {
	if upstream.ID == "" {
		return errors.New("upstream ID不能为空")
	}
	global.Upstreams.Store(upstream.ID, upstream)

	// 更新中间件
	mwLen := len(upstream.Middleware)
	if mwLen == 0 {
		global.UpstreamMiddleware.Delete(upstream.ID)
		return nil
	}
	mw := make([]global.MiddlewareType, mwLen)
	for k := range upstream.Middleware {
		m, err := middleware.Build(upstream.Middleware[k].Name, upstream.Middleware[k].Config, false)
		if err != nil {
			return err
		}
		mw = append(mw, m)
	}
	global.UpstreamMiddleware.Store(upstream.ID, mw)
	return nil
}

func DelUpstream(upstreamID string) error {
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
