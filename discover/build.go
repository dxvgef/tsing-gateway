package discover

import (
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/discover/tsing_center"
	"github.com/dxvgef/tsing-gateway/global"
)

// 构建节点发现实例
// key为节点发现方式的名称，value为节点发现的参数json字符串
func Build(name, config string) (global.DiscoverType, error) {
	switch name {
	case "tsing_center":
		f, err := tsing_center.New(config)
		if err != nil {
			log.Err(err).Caller().Send()
			return nil, err
		}
		return f, nil
	}
	return nil, errors.New("没有找到名为" + name + "的探测器")
}
