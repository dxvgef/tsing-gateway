package favicon

import (
	"encoding/json"
	"log"
	"net/http"
)

// 收藏夹图标处理
type Favicon struct {
	ReCode int `json:"re_code"`
}

// 获得中间件实例
func Inst(config string) (*Favicon, error) {
	var filter Favicon
	err := json.Unmarshal([]byte(config), &filter)
	if err != nil {
		return nil, err
	}
	return &filter, nil
}

// 中间件行为
func (f *Favicon) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	if req.RequestURI == "/favicon.ico" {
		log.Println("favicon中间件触发了")
		resp.WriteHeader(f.ReCode)
		return false, nil
	}
	return true, nil
}
