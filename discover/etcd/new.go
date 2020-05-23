package etcd

import (
	"encoding/json"
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
