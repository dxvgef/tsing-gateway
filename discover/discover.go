package discover

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover/coredns_etcd"
)

// 节点发现接口
type Discover interface {
	Action() (string, int, int, int, error)
}

// 构建节点发现实例
// key为节点发现方式的名称，value为节点发现的参数json字符串
func Build(name, config string) (Discover, error) {
	switch name {
	case "coredns_etcd":
		f, err := coredns_etcd.New(config)
		if err != nil {
			log.Error().Caller().Msg(err.Error())
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("not found endpoint by name")
}
