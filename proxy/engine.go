package proxy

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover"
	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/load_balance"
)

// 代理引擎
type Engine struct{}

// 实现http.Handler接口的方法
// 下游请求入口
func (*Engine) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	var (
		next bool
		err  error
	)
	// 匹配到上游
	upstream, status := matchRoute(req)
	if status > 0 {
		resp.WriteHeader(status)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(status))); err != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 执行全局中间件
	for k := range global.GlobalMiddleware {
		next = false
		next, err = global.GlobalMiddleware[k].Action(resp, req)
		if err != nil {
			log.Err(err).Caller().Str("name", global.GlobalMiddleware[k].GetName()).Msg("执行全局中间件时出错")
			return
		}
		if !next {
			return
		}
	}

	// 执行上游中间件
	// k := 0
	global.UpstreamMiddleware.Range(func(k, v interface{}) bool {
		if k.(string) != upstream.ID {
			return true
		}
		mw := v.(global.MiddlewareType)
		next = false
		// 执行中间件逻辑
		next, err = mw.Action(resp, req)
		if err != nil {
			log.Err(err).Caller().Str("upstream id", upstream.ID).Str("middleware name", mw.GetName()).Msg("执行上游中间件时出错")
			return false
		}
		if !next {
			return false
		}
		return true
	})
	if err != nil {
		return
	}

	// 获得端点
	var endpointURL *url.URL
	endpointURL, status, err = getEndpointURL(upstream)
	if err != nil {
		log.Err(err).Caller().Msg("获得端点失败")
		resp.WriteHeader(status)
		if _, err = resp.Write(global.StrToBytes(http.StatusText(status))); err != nil {
			log.Err(err).Caller().Send()
		}
		return
	}

	// 转发请求到端点
	p := httputil.NewSingleHostReverseProxy(endpointURL)
	req.URL.Host = endpointURL.Host
	req.URL.Scheme = endpointURL.Scheme
	p.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		log.Debug().Err(err).Caller().Msg("向端点发起请求出错")
	}
	p.ServeHTTP(resp, req)
}

// 获得端点url
func getEndpointURL(upstream global.UpstreamType) (endpointURL *url.URL, status int, err error) {
	// 获得所有端点列表
	var endpoints []global.EndpointType
	endpoints, status, err = fetchEndpoints(upstream)
	if status > 0 {
		return
	}

	if len(endpoints) == 1 {
		endpointURL, err = url.Parse(endpoints[0].Addr)
		if err != nil {
			status = http.StatusInternalServerError
			log.Err(err).Str("addr", endpoints[0].Addr).Caller().Msg("解析端点地址失败")
			return
		}
		return
	}

	var lb global.LoadBalance
	// 使用负载均衡算法选取一个最终的端点
	lb, err = load_balance.Use(upstream.LoadBalance)
	if err != nil {
		status = http.StatusServiceUnavailable
		log.Err(err).Str("upstream id", upstream.ID).Str("load_balance", upstream.LoadBalance).Caller().Msg("使用负载均衡算法失败")
		return
	}

	for i := range endpoints {
		lb.Put(upstream.ID, endpoints[i].Addr, endpoints[i].Weight)
	}
	endpointAddr := lb.Next(upstream.ID)
	endpointURL, err = url.Parse(endpointAddr)
	if err != nil {
		status = http.StatusInternalServerError
		log.Err(err).Str("addr", endpointAddr).Caller().Msg("解析端点地址失败")
		return
	}

	log.Debug().Caller().Str("addr", endpointURL.String()).Msg("目标端点")
	return
}

// 从上游获得所有端点
func fetchEndpoints(upstream global.UpstreamType) (endpoints []global.EndpointType, status int, err error) {
	// 如果有静态端点则静态优先
	if upstream.StaticEndpoint != "" {
		endpoints = append(endpoints, global.EndpointType{
			UpstreamID: upstream.ID,
			Addr:       upstream.StaticEndpoint,
			Weight:     0,
		})
		return
	}

	if upstream.Discover.Name == "" {
		status = http.StatusInternalServerError
		err = errors.New("上游没有设置静态端点也没有配置探测器")
		log.Err(err).Caller().Str("upstream id", upstream.ID).Str("config", upstream.Discover.Config).Msg("构建探测器时出错")
		return
	}
	// 构建探测器
	var dc global.DiscoverType
	dc, err = discover.Build(upstream.Discover.Name, upstream.Discover.Config)
	if err != nil {
		status = http.StatusInternalServerError
		log.Err(err).Caller().Str("name", upstream.Discover.Name).Str("config", upstream.Discover.Config).Msg("构建探测器时出错")
		return
	}
	// 获得所有端点
	endpoints, err = dc.Fetch(upstream.ID)
	if err != nil {
		status = http.StatusInternalServerError
		log.Err(err).Caller().Str("name", upstream.Discover.Name).Str("config", upstream.Discover.Config).Msg("获得端点列表时出错")
	}
	return
}
