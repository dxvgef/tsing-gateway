package global

import (
	"net/http"
)

// 用于构建服务、中间件、存储器模块时的参数配置
type ModuleConfig struct {
	Name   string `json:"name,omitempty"`
	Config string `json:"config,omitempty"`
}

// 主机
type HostType struct {
	RouteGroupID string         `json:"route_group_id"`       // 路由分组ID
	Middleware   []ModuleConfig `json:"middleware,omitempty"` // 中间件配置
}

// 服务
type ServiceType struct {
	ID             string         `json:"id"`                        // 服务ID
	Middleware     []ModuleConfig `json:"middleware,omitempty"`      // 中间件配置
	StaticEndpoint string         `json:"static_endpoint,omitempty"` // 静态端点地址，优先级高于Discover
	Discover       ModuleConfig   `json:"discover,omitempty"`        // 探测器配置
}

// 端点发现
type DiscoverType interface {
	Fetch(string) (NodeType, error)
}

// 中间件接口
type MiddlewareType interface {
	Action(http.ResponseWriter, *http.Request) (bool, error)
	GetName() string
}

type NodeType struct {
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
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

	LoadAllService() error             // 从存储器加载所有服务到本地
	LoadService([]byte) error          // 从存储器加载单个服务数据
	SaveAllService() error             // 将本地所有服务数据保存到存储器
	SaveService(string, string) error  // 将本地单个服务数据保存到存储器
	DeleteLocalService(string) error   // 删除本地单个服务数据
	DeleteStorageService(string) error // 删除存储器中单个服务数据

	LoadAllRoute() error                             // 从存储器加载所有路由数据到本地
	LoadRoute(string, []byte) error                  // 从存储器加载单个路由数据
	SaveAllRoute() error                             // 将本地所有路由保存到存储器
	SaveRoute(string, string, string, string) error  // 将本地单个路由数据保存到存储器
	DeleteLocalRoute(string) error                   // 删除本地单个路由数据
	DeleteStorageRoute(string, string, string) error // 删除存储器中单个路由数据

	Watch() error // 监听存储器的数据变更
}
