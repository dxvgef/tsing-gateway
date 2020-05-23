package health

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

// 健康检查
type Health struct {
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

// 新建中间件实例
func New(config string) (*Health, error) {
	var instance Health
	err := json.Unmarshal([]byte(config), &instance)
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (self *Health) GetName() string {
	return "health"
}

func (self *Health) GetConfig() ([]byte, error) {
	return self.MarshalJSON()
}

// 中间件行为
func (self *Health) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	log.Debug().Msg("执行了中间件：health")
	return true, nil
}
