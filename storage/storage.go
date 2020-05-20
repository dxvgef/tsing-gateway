package storage

import (
	"errors"

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
			return nil, err
		}
		global.StorageKeyPrefix = sa.KeyPrefix
		return sa, nil
	}
	return nil, errors.New("根据名称没有找到对应的存储器")
}
