package global

import "net/http"

type UpstreamType struct {
	ID         string         `json:"id"`                   // 上游ID
	Middleware []ModuleConfig `json:"middleware,omitempty"` // 中间件配置
	Discover   ModuleConfig   `json:"discover"`             // 节点发现配置
	// 启用缓存，如果关闭，则每次请求都从etcd中获取endpoints
	Cache bool `json:"cache"`
	/*
		缓存重试次数
		在缓存中失败达到指定次数后，重新从discover中获取endpoints来更新缓存
	*/
	CacheRetry   int            `json:"cache_retry"`
	Endpoints    []EndpointType `json:"-"`                      // 终点列表
	LoadBalance  string         `json:"load_balance,omitempty"` // 负载均衡算法
	LastEndpoint string         `json:"-"`                      // 最后使用的endpoint，用于防止连续命中同一个
}

// 模块配置，用于upstream
type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}

// 终点
type EndpointType struct {
	UpstreamID string `json:"upstream_id"`
	URL        string `json:"url"`
	Weight     int    `json:"weight"`
}

// 节点发现接口
type DiscoverType interface {
	Fetch() (EndpointType, error)
	FetchAll() ([]EndpointType, error)
}

// 定义中间件接口
type MiddlewareType interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
}

// 存储器接口
type StorageType interface {
	LoadAll() error                                // 加载所有数据
	LoadAllMiddleware() error                      // 加载所有middlware数据
	LoadAllUpstreams() error                       // 加载所有upstream数据
	LoadAllRoutes() error                          // 加载所有route数据
	LoadAllHosts() error                           // 加载所有host数据
	SaveAll() error                                // 存储所有数据
	SaveAllUpstreams() error                       // 存储所有upstream数据
	SaveAllRoutes() error                          // 存储所有route数据
	SaveAllHosts() error                           // 存储所有host数据
	Watch() error                                  // 监听数据变更
	PutHost(string, string) error                  // 设置单个主机
	DelHost(string) error                          // 删除单个主机
	PutUpstream(string, string) error              // 设置单个上游
	DelUpstream(string) error                      // 删除单个上游
	PutRoute(string, string, string, string) error // 设置单个路由
	DelRoute(string, string, string) error         // 删除单个路嵋
}
