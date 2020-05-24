package global

import (
	"github.com/bwmarrin/snowflake"
)

var (
	SnowflakeNode *snowflake.Node

	Storage          StorageType // 存储器
	StorageKeyPrefix string      // 存储器键名前缀
	// 存储器客户端ID，在存储器被构建时自动生成
	// 用于存储器的
	StorageClientID    int64
	GlobalMiddleware   []MiddlewareType                                // 全局中间件
	Hosts              = make(map[string]string)                       // 主机列表
	Routes             = make(map[string]map[string]map[string]string) // 路由列表
	Upstreams          = make(map[string]UpstreamType)                 // 上游列表
	UpstreamMiddleware = make(map[string][]MiddlewareType)             // 所有上游的中间件实例
	// Endpoints          = make(map[string][]EndpointType)               // 所有端点列表
	// HTTP方法允许的值
	Methods = []string{
		"ANY", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
	}
	// 负载均衡算数允许的值
	// LoadBalance = []string{"discover", "wred"}
)
