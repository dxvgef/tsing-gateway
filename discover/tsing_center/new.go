package tsing_center

import (
	"encoding/json"
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
		return nil, err
	}
	return &e, nil
}
