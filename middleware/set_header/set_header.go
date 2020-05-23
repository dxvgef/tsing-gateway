package set_header

import (
	"net/http"

	"github.com/dxvgef/tsing-gateway/global"
)

// header数据处理
type SetHeader struct {
	Request  map[string]string `json:"request,omitempty"`
	Response map[string]string `json:"response,omitempty"`
}

// 新建中间件实例
func New(config string) (*SetHeader, error) {
	var instance SetHeader
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

func (self *SetHeader) GetName() string {
	return "set_header"
}

// 中间件行为
func (self *SetHeader) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	for k, v := range self.Request {
		req.Header.Set(k, v)
	}
	for k, v := range self.Response {
		resp.Header().Set(k, v)
	}
	return true, nil
}
