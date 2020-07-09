package cors

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"local/global"
)

// CORS
type CORS struct {
	AllowOrigins     []string `json:"allow_origins,omitempty"`     // 允许客户端的来源域
	ExposeHeaders    []string `json:"expose_headers,omitempty"`    // 允许响应的头信息
	AllowCredentials bool     `json:"allow_credentials,omitempty"` // 允许客户端携带cookie
	AllowMethods     []string `json:"allow_methods,omitempty"`     // 允许客户端请求的方法
	AllowHeaders     []string `json:"allow_headers,omitempty"`     // 允许客户端请求的头信息
}

// 新建中间件实例
func New(config string) (*CORS, error) {
	var instance CORS
	instance.AllowOrigins = []string{"*"}
	instance.AllowMethods = []string{"*"}
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
	if len(self.AllowOrigins) > 0 {
		if self.AllowMethods[0] == "*" {
			resp.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			resp.Header().Set("Access-Control-Allow-Origin", strings.Join(self.AllowOrigins, ","))
		}
	}
	if len(self.AllowMethods) > 0 {
		if self.AllowMethods[0] == "*" {
			resp.Header().Set("Access-Control-Allow-Methods", "*")
		} else {
			resp.Header().Set("Access-Control-Allow-Methods", strings.Join(self.AllowMethods, ","))
		}
	}
	if len(self.AllowHeaders) > 0 {
		if self.AllowHeaders[0] == "*" {
			resp.Header().Set("Access-Control-Allow-Headers", "*")
		} else {
			resp.Header().Set("Access-Control-Allow-Headers", strings.Join(self.AllowHeaders, ","))
		}
	}
	if len(self.ExposeHeaders) > 0 {
		if self.ExposeHeaders[0] == "*" {
			resp.Header().Set("Access-Control-Expose-Headers", "*")
		} else {
			resp.Header().Set("Access-Control-Expose-Headers", strings.Join(self.ExposeHeaders, ","))
		}
	}
	if len(self.ExposeHeaders) > 0 {
		if self.ExposeHeaders[0] == "*" {
			resp.Header().Set("Access-Control-Expose-Headers", "*")
		} else {
			resp.Header().Set("Access-Control-Expose-Headers", strings.Join(self.ExposeHeaders, ","))
		}
	}
	if self.AllowCredentials {
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	return false, nil
}
