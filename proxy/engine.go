package proxy

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
)

// 代理引擎
type Engine struct{}

// 实现http.Handler接口的方法
// 下游请求入口
func (*Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	// 通过路由匹配到上游
	upstream, status := matchRoute(req)
	if status != http.StatusOK {
		resp.WriteHeader(status)
		// nolint
		_, _ = resp.Write(global.StrToBytes(http.StatusText(status)))
		return
	}

	// 执行全局中间件
	for k := range global.GlobalMiddleware {
		next, err := global.GlobalMiddleware[k].Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// 执行上游中间件
	for k := range global.UpstreamMiddleware[upstream.ID] {
		// 执行中间件逻辑
		next, err := global.UpstreamMiddleware[upstream.ID][k].Action(resp, req)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return
		}
		if !next {
			return
		}
	}

	// 根据上游中的探测器配置实时构建探测器实例
	_, err := discover.Build(upstream.Discover.Name, upstream.Discover.Config)
	if err != nil {
		log.Error().Caller().Msg(err.Error())
		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write(global.StrToBytes(http.StatusText(http.StatusInternalServerError))) // nolint
		return
	}

	// todo 这里要写获取endpoint的逻辑

	// todo 以下是反向代理的请求逻辑，暂时用200状态码替代
	resp.WriteHeader(http.StatusOK)
	// if _, err := resp.Write(global.StrToBytes(respText)); err != nil {
	// 	log.Error().Msg(err.Error())
	// }

	// 这里使用的servHTTP是一个使用新协程的非阻塞处理方式
	// resp.Header().Update("X-Power-By", "Tsing Gateway")
	// p.ServeHTTP(resp, req)
}
