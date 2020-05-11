package explorer

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/explorer/coredns_etcd"
)

type Explorer interface {
	Action() (string, int, int, int, error)
}

// 构建探测器实例
// key为探测器的名称，value为探测器的参数json字符串
func Build(name, config string) (Explorer, error) {
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
