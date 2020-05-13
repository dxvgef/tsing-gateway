package source

import (
	"errors"

	"github.com/dxvgef/tsing-gateway/proxy"
	"github.com/dxvgef/tsing-gateway/source/etcd"
)

// 数据源接口
type Source interface {
	LoadAll() error          // 加载所有数据
	LoadAllUpstreams() error // 加载所有upstream数据
	LoadAllRoutes() error    // 加载所有route数据
	LoadAllHosts() error     // 加载所有host数据
	SaveAll() error          // 存储所有数据
	SaveAllUpstreams() error // 存储所有upstream数据
	SaveAllRoutes() error    // 存储所有route数据
	SaveAllHosts() error     // 存储所有host数据
}

// 构建数据源实例
// key为数据源的名称，value为数据源的参数json字符串
func Build(e *proxy.Engine, name, config string) (Source, error) {
	switch name {
	case "etcd":
		f, err := etcd.New(e, config)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("not found source by name")
}
