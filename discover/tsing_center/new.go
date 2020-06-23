package tsing_center

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

// tsing center
type TsingCenter struct {
	Addr   string `json:"addr"`   // tsing center的地址
	Secret string `json:"secret"` // tsing center的连接密钥
}

// 新建探测器实例
func New(config string) (*TsingCenter, error) {
	var e TsingCenter
	err := json.Unmarshal([]byte(config), &e)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &e, nil
}
