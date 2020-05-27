package etcd

import (
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/rs/zerolog/log"
)

// etcd
type Etcd struct {
	EtcdEndpoints []string         `json:"etcd_endpoints"` // etcd的endpoints
	KeyPrefix     string           `json:"key_prefix"`     // 键名前缀
	client        *clientv3.Client // etcd客户端
}

// 新建探测器实例
func New(config string) (*Etcd, error) {
	var e Etcd
	err := json.Unmarshal([]byte(config), &e)
	if err != nil {
		return nil, err
	}

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints: e.EtcdEndpoints,
	})
	if err != nil {
		log.Err(err).Caller().Msg("etcd连接失败")
		return nil, err
	}
	return &e, nil
}
