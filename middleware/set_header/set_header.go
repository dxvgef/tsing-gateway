package set_header

import (
	"net/http"

	"github.com/rs/zerolog/log"

	"local/global"
)

// header数据处理
type SetHeader struct {
	RequestHeader  map[string]string `json:"request_header,omitempty"`
	ResponseHeader map[string]string `json:"response_header,omitempty"`
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
	log.Debug().Caller().Msg("set_header")
	for k, v := range self.RequestHeader {
		req.Header.Set(k, v)
	}
	for k, v := range self.ResponseHeader {
		resp.Header().Set(k, v)
	}
	return
}
