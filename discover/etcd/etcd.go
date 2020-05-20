package etcd

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/discover"
)

// etcd
type Etcd struct {
	Hosts     []string `json:"hosts"`      // etcd的endpoints
	KeyPrefix string   `json:"key_prefix"` // 键名前缀
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
func (self *Etcd) Fetch() (endpoint discover.Endpoint, err error) {
	return
}

// 获取所有Endpoint
func (self *Etcd) FetchAll() (endpoints []discover.Endpoint, err error) {
	return
}
