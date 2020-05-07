package header

import (
	"encoding/json"
	"net/http"
)

// header数据处理
type Header struct {
	Request  map[string]string `json:"request,omitempty"`
	Response map[string]string `json:"response,omitempty"`
}

// 获得中间件实例
func Inst(config string) (*Header, error) {
	var mw Header
	err := json.Unmarshal([]byte(config), &mw)
	if err != nil {
		return nil, err
	}
	return &mw, nil
}

// 中间件行为
func (mw *Header) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	for k, v := range mw.Request {
		req.Header.Set(k, v)
	}
	for k, v := range mw.Response {
		resp.Header().Set(k, v)
	}
	return true, nil
}
