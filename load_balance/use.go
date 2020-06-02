package load_balance

import (
	"errors"
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/load_balance/swrr"
	"github.com/dxvgef/tsing-gateway/load_balance/wr"
	"github.com/dxvgef/tsing-gateway/load_balance/wrr"
)

// 使用指定算法的负载均衡
func Use(name string) (global.LoadBalance, error) {
	name = strings.ToUpper(name)
	switch name {
	// 加权随机
	case "WR":
		return wr.Init(), nil
	// 加权轮循
	case "WRR":
		return wrr.Init(), nil
	// 平滑加权轮循
	case "SWRR":
		return swrr.Init(), nil
	}
	return nil, errors.New("不支持的负载均衡算法")
}
