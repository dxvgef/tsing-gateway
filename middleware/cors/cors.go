package cors

import (
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"local/global"
)

// CORS
type CORS struct {
	AllowOrigins     string `json:"allow_origins,omitempty"`     // 允许客户端的来源域
	ExposeHeaders    string `json:"expose_headers,omitempty"`    // 允许响应的头信息
	AllowCredentials bool   `json:"allow_credentials,omitempty"` // 允许客户端携带cookie
	AllowMethods     string `json:"allow_methods,omitempty"`     // 允许客户端请求的方法
	AllowHeaders     string `json:"allow_headers,omitempty"`     // 允许客户端请求的头信息
}

// 新建中间件实例
func New(config string) (*CORS, error) {
	var instance CORS
	instance.AllowOrigins = "*"
	instance.AllowMethods = "GET,POST,PUT,DELETE,OPTIONS,PATCH"
	instance.AllowHeaders = "Access-Control-Allow-Headers,token,Origin,X-Requested-With,Content-Type,Accept,X-Token"
	instance.AllowCredentials = true
	instance.ExposeHeaders = "*"
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}

func (self *CORS) GetName() string {
	return "cors"
}

// 中间件行为
func (self *CORS) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	resp.Header().Set("Access-Control-Allow-Origin", self.AllowOrigins)
	resp.Header().Set("Access-Control-Allow-Methods", self.AllowMethods)
	resp.Header().Set("Access-Control-Allow-Headers", self.AllowHeaders)
	resp.Header().Set("Access-Control-Expose-Headers", self.ExposeHeaders)
	resp.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(self.AllowCredentials))
	if req.Method == "OPTIONS" {
		resp.WriteHeader(http.StatusNoContent)
		return true, nil
	}
	return false, nil
}
