package load_balance

import (
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/load_balance/swrr"
	"github.com/dxvgef/tsing-gateway/load_balance/wr"
	"github.com/dxvgef/tsing-gateway/load_balance/wrr"
)

// 使用指定算法的负载均衡
func Use(name string) global.LoadBalance {
	name = strings.ToUpper(name)
	switch name {
	// 加权随机
	case "WR":
		return wr.Init()
	// 加权轮循
	case "WRR":
		return wrr.Init()
	// 平滑加权轮循
	case "SWRR":
		return swrr.Init()
	}
	return nil
}
