package health_check

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// 健康检查
type HealthCheck struct {
	Active struct {
		On               bool   `json:"on,omitempty"`       // 打开健康检查
		Interval         int    `json:"interval,omitempty"` // 检查间隔的时间(秒)
		URL              string `json:"url,omitempty"`      // 主动检查地址
		FailureStateCode []int  `json:"failure_state_code,omitempty"`
	} `json:"active,omitempty"`
	Passive struct {
		On  bool `json:"on,omitempty"`  // 打开健康检查
		TTL int  `json:"ttl,omitempty"` // 端点的生命周期(秒)
	} `json:"passive,omitempty"`
}

// 获得中间件实例
func Inst(config string) (*HealthCheck, error) {
	var filter HealthCheck
	err := json.Unmarshal([]byte(config), &filter)
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// 中间件行为
func (f *HealthCheck) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	log.Debug().Msg("执行了过滤器：health_check")
	return true, nil
}
