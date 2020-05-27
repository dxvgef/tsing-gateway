package coredns_etcd

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/global"
)

// coredns etcd
type CoreDNSEtcd struct {
	EtcdEndpoints []string `json:"etcd_endpoints"`
	KeyPrefix     string   `json:"key_prefix"`
	Host          string   `json:"host"`
}

// 新建探测器实例
func New(config string) (*CoreDNSEtcd, error) {
	var e CoreDNSEtcd
	err := json.Unmarshal([]byte(config), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// 获取单个endpoint
func (self *CoreDNSEtcd) Fetch(upstreamID string) (endpoint []global.EndpointType, err error) {
	return
}
