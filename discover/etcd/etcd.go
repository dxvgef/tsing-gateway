package etcd

import (
	"encoding/json"
)

// etcd
type Etcd struct {
	EtcdEndpoints []string `json:"etcd_endpoints"`
	KeyPrefix     string   `json:"key_prefix"`
	Host          string   `json:"host"`
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

// 探测行为
func (self *Etcd) Action() (ip string, port int, weight int, ttl int, err error) {
	return
}
