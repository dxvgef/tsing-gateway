package global

import "github.com/bwmarrin/snowflake"

// 节点ID
var ID int64

// ID生成器的实例
var IDNode *snowflake.Node

// 存储器内的键名前缀
var StorageKeyPrefix string

// http方法
var Methods = []string{
	"*", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
}

// 负载均衡算法
var LoadBalance = []string{"discover", "wred"}

// 模块配置，用于upstream
type ModuleConfig struct {
	Name   string `json:"name"`
	Config string `json:"config"`
}
type Endpoint struct {
	UpstreamID string `json:"upstream_id"`
	URL        string `json:"url"`
	Weight     int    `json:"weight"`
}
