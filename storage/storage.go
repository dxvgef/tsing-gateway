package storage

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/storage/etcd"
)

// 构建存储器实例
// key为存储器的名称，value为存储器的参数json字符串
func Build(name, config string) (global.StorageType, error) {
	switch name {
	case "etcd":
		sa, err := etcd.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		return sa, nil
	}
	return nil, errors.New("不支持的存储器名称")
}
