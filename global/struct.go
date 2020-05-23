package global

import "net/http"

// 用于构建上游、中间件、存储器模块时的参数配置
type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// 上游
type UpstreamType struct {
	ID             string         `json:"id"`                        // 上游ID
	Middleware     []ModuleConfig `json:"middleware,omitempty"`      // 中间件配置
	StaticEndpoint string         `json:"static_endpoint,omitempty"` // 静态端点地址，如果有值，则不使用探测器发现节点
	Discover       ModuleConfig   `json:"discover,omitempty"`        // 探测器配置
	// 启用缓存，如果关闭，则每次请求都从etcd中获取端点
	// Cache bool `json:"cache"`
	/*
		缓存重试次数
		在缓存中失败达到指定次数后，重新从discover中获取endpoints来更新缓存
	*/
	// CacheRetry   int            `json:"cache_retry"`
	// Endpoints []EndpointType `json:"-"` // 终点列表
	// LoadBalance  string         `json:"load_balance,omitempty"` // 负载均衡算法
	// LastEndpoint string         `json:"-"`                      // 最后使用的端点，用于防止连续命中同一个
}

// 端点
type EndpointType struct {
	UpstreamID string `json:"upstream_id"`
	URL        string `json:"url"`
	Weight     int    `json:"weight"`
}

// 端点发现
type DiscoverType interface {
	Fetch(string) (EndpointType, error)
	FetchAll(string) ([]EndpointType, error)
}

// 中间件接口
type MiddlewareType interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
	GetName() string
}

// 存储器
type StorageType interface {
	LoadAll() error          // 加载所有数据
	LoadAllUpstreams() error // 加载所有上游数据
	LoadAllRoutes() error    // 加载所有路由数据
	LoadAllHosts() error     // 加载所有主机数据
	LoadMiddleware() error   // 加载全局中间件数据

	SaveAll() error          // 存储所有数据
	SaveAllUpstreams() error // 存储所有上游数据
	SaveAllRoutes() error    // 存储所有路由数据
	SaveAllHosts() error     // 存储所有主机数据
	SaveMiddleware() error   // 存储全局中间件

	Watch() error // 监听数据变更

	PutHost(string, string) error // 设置单个主机
	DelHost(string) error         // 删除单个主机

	PutUpstream(string, string) error // 设置单个上游
	DelUpstream(string) error         // 删除单个上游

	PutRoute(string, string, string, string) error // 设置单个路由
	DelRoute(string, string, string) error         // 删除路由

	PutMiddleware(string) error // 设置全局中间件
}
