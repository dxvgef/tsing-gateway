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
	RouteGroupID string         `json:"route_group_id"`       // 路由分组ID
	Middleware   []ModuleConfig `json:"middleware,omitempty"` // 中间件配置
}

// 上游
type UpstreamType struct {
	ID             string         `json:"id"`                        // 上游ID
	Middleware     []ModuleConfig `json:"middleware,omitempty"`      // 中间件配置
	StaticEndpoint string         `json:"static_endpoint,omitempty"` // 静态端点地址，优先级高于Discover
	Discover       ModuleConfig   `json:"discover,omitempty"`        // 探测器配置
	LoadBalance    string         `json:"load_balance,omitempty"`    // 负载均衡算法名称

	/*
		最大缓存容错次数
		缓存中的端点不可用数超过此值，则自动从上游重新获取端点列表来更新希尔顿
	*/
	MaxCacheFault int `json:"max_cache_fault"`
}

// 负载均衡接口
type LoadBalance interface {
	Set(string, string, int) error
	Remove(string, string) error
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
	LoadAll() error // 从存储器加载所有数据到本地
	SaveAll() error // 将本地所有数据保存到存储器

	LoadAllHost() error             // 从存储器加载所有主机数据到本地
	SaveAllHost() error             // 将本地所有主机数据保存到存储器
	LoadHost(string, []byte) error  // 从存储器加载单个主机数据
	SaveHost(string, string) error  // 将本地单个主机数据保存到存储器
	DeleteLocalHost(string) error   // 删除本地单个主机数据
	DeleteStorageHost(string) error // 删除存储器中单个主机数据

	LoadAllUpstream() error             // 从存储器加载所有上游到本地
	LoadUpstream([]byte) error          // 从存储器加载单个上游数据
	SaveAllUpstream() error             // 将本地所有上游数据保存到存储器
	SaveUpstream(string, string) error  // 将本地单个上游数据保存到存储器
	DeleteLocalUpstream(string) error   // 删除本地单个上游数据
	DeleteStorageUpstream(string) error // 删除存储器中单个上游数据

	LoadAllRoute() error                             // 从存储器加载所有路由数据到本地
	LoadRoute(string, []byte) error                  // 从存储器加载单个路由数据
	SaveAllRoute() error                             // 将本地所有路由保存到存储器
	SaveRoute(string, string, string, string) error  // 将本地单个路由数据保存到存储器
	DeleteLocalRoute(string) error                   // 删除本地单个路由数据
	DeleteStorageRoute(string, string, string) error // 删除存储器中单个路由数据

	Watch() error // 监听存储器的数据变更
}
