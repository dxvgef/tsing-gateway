package tsing_center

import (
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

// 获节点
func (self *TsingCenter) Fetch(serviceID string) (global.NodeType, error) {
	var node global.NodeType
	endpoint := self.Addr + "/services/" + serviceID + "/select"
	resp, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Err(err).Caller().Msg("探测器发起请求失败")
		return node, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Caller().Msg("解析探测器响应数据失败")
		return node, err
	}
	err = node.UnmarshalJSON(body)
	if err != nil {
		log.Err(err).Caller().Msg("解码探测器响应数据失败")
		return node, err
	}
	return node, nil
}
