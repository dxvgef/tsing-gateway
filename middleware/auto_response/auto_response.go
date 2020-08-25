package auto_response

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"local/global"
)

type Rule struct {
	Method string `json:"method"`           // 触发自动响应的请求方法，大写，允许*匹配所有
	Status int    `json:"status,omitempty"` // 自动响应的状态码
	Data   string `json:"data,omitempty"`   // 自动响应的内容
}

type AutoResponse struct {
	data map[string]Rule // key为路径，允许*匹配所有
}

func New(config string) (*AutoResponse, error) {
	var instance AutoResponse
	err := json.Unmarshal(global.StrToBytes(config), &instance.data)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}

func (self *AutoResponse) GetName() string {
	return "auto_response"
}

func (self *AutoResponse) Action(resp http.ResponseWriter, req *http.Request) (abort bool, err error) {
	for k := range self.data {
		if (k != "*" && req.RequestURI != k) || (self.data[k].Method != "ANY" && req.Method != self.data[k].Method) {
			continue
		}
		if self.data[k].Status != 0 {
			resp.WriteHeader(self.data[k].Status)
			abort = true
		}
		if self.data[k].Data != "" {
			abort = true
			_, err = resp.Write(global.StrToBytes(self.data[k].Data))
			if err != nil {
				log.Err(err).Caller().Send()
			}
		}
		return
	}
	return
}
