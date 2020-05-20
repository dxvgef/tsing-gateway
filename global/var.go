package global

import (
	"github.com/bwmarrin/snowflake"
)

var (
	ID               int64
	StorageKeyPrefix string
	IDNode           *snowflake.Node
	Methods          = []string{
		"*", "GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "TRACE", "CONNECT",
	}

	Storage StorageType

	LoadBalance = []string{"discover", "wred"}

	Middleware []ModuleConfig

	Hosts map[string]string

	Routes map[string]map[string]map[string]string

	Upstreams map[string]UpstreamType
)
