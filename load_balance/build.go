package load_balance

import (
	"strings"

	"github.com/dxvgef/tsing-gateway/global"
	"github.com/dxvgef/tsing-gateway/load_balance/swrr"
	"github.com/dxvgef/tsing-gateway/load_balance/wr"
	"github.com/dxvgef/tsing-gateway/load_balance/wrr"
)

func Build(name string) global.LoadBalance {
	name = strings.ToUpper(name)
	switch name {
	// 加权随机
	case "WR":
		return wr.New()
	// 加权轮循
	case "WRR":
		return wrr.New()
	// 平滑加权轮循
	case "SWRR":
		return swrr.New()
	}
	return nil
}
