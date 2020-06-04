package global

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	SnowflakeNode *snowflake.Node

	Storage StorageType // 存储器

	Hosts              sync.Map // 主机 key=hostname, value=HostType
	HostMiddleware     sync.Map // 主机中间件 key=hostname, value=[]MiddlewareType
	Routes             sync.Map // 路由 key=routeGroupID@path@method, value=upstreamID
	Upstreams          sync.Map // 上游 key=RouteGroupID, value=UpstreamType
	UpstreamMiddleware sync.Map // 上游中间件 key=MiddlewareName, value=[]MiddlewareType

	// HTTP方法允许的值
	HTTPMethods = []string{
		"ANY", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
	}
)
