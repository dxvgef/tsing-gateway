package set_header

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"local/global"
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
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}

func (self *SetHeader) GetName() string {
	return "set_header"
}

// 中间件行为
func (self *SetHeader) Action(resp http.ResponseWriter, req *http.Request) (abort bool, err error) {
	for k, v := range self.Request {
		req.Header.Set(k, v)
	}
	for k, v := range self.Response {
		resp.Header().Set(k, v)
	}
	return
}
