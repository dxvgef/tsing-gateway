package global

import (
	"net/http"
)

// 用于构建上游、中间件、存储器模块时的参数配置
type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// 主机
type HostType struct {
	UpstreamID string         `json:"upstream_id"`          // 上游ID
	Middleware []ModuleConfig `json:"middleware,omitempty"` // 中间件配置
}

// 上游
type UpstreamType struct {
	ID             string         `json:"id"`                        // 上游ID
	Middleware     []ModuleConfig `json:"middleware,omitempty"`      // 中间件配置
	StaticEndpoint string         `json:"static_endpoint,omitempty"` // 静态端点地址，优先级高于Discover
	Discover       ModuleConfig   `json:"discover,omitempty"`        // 探测器配置
	LoadBalance    string         `json:"load_balance"`              // 负载均衡算法名称

	/*
		最大缓存容错次数
		缓存中的端点不可用数超过此值，则自动从上游重新获取端点列表来更新希尔顿
	*/
	MaxCacheFault int `json:"max_cache_fault"`
}

// 负载均衡接口
type LoadBalance interface {
	Add(string, string, int) error
	Put(string, string, int)
	Next(string) string
	Total(string) int
}

// 端点
type EndpointType struct {
	UpstreamID string `json:"upstream_id"`
	Addr       string `json:"addr"`
	Weight     int    `json:"weight"`
}

// 端点发现
type DiscoverType interface {
	Fetch(string) ([]EndpointType, error)
}

// 中间件接口
type MiddlewareType interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
	GetName() string
}

// 存储器
type StorageType interface {
	LoadAll() error          // 加载所有数据
	LoadAllHosts() error     // 加载所有主机数据
	LoadAllUpstreams() error // 加载所有上游数据
	LoadAllRoutes() error    // 加载所有路由数据

	SaveAll() error          // 存储所有数据
	SaveAllUpstreams() error // 存储所有上游数据
	SaveAllRoutes() error    // 存储所有路由数据
	SaveAllHosts() error     // 存储所有主机数据

	PutHost(string, string) error // 设置单个主机
	DelHost(string) error         // 删除单个主机

	PutUpstream(string, string) error // 设置单个上游
	DelUpstream(string) error         // 删除单个上游

	PutRoute(string, string, string, string) error // 设置单个路由
	DelRoute(string, string, string) error         // 删除单个路由

	Watch() error // 监听数据变更
}
