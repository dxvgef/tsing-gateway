package etcd

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/global"
)

// etcd
type Etcd struct {
	EtcdEndpoints []string `json:"etcd_endpoints"` // etcd的endpoints
	KeyPrefix     string   `json:"key_prefix"`     // 键名前缀
}

// 新建探测器实例
func New(config string) (*Etcd, error) {
	var e Etcd
	err := json.Unmarshal([]byte(config), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// 获取单个Endpoint
func (self *Etcd) Fetch() (endpoint global.EndpointType, err error) {
	return
}

// 获取所有Endpoint
func (self *Etcd) FetchAll() (endpoints []global.EndpointType, err error) {
	return
}
