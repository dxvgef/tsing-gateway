package cors

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"local/global"
)

// CORS
type CORS struct {
	AllowOrigins     []string `json:"allow_origins,omitempty"` // 允许的来源
	ExposeHeaders    []string `json:"expose_headers,omitempty"`
	AllowCredentials bool     `json:"allow_credentials,omitempty"` // 允许携带cookie
	AllowMethods     []string `json:"allow_methods,omitempty"`     // 允许的方法
	AllowHeaders     []string `json:"allow_headers,omitempty"`     // 允许的header
}

// 新建中间件实例
func New(config string) (*CORS, error) {
	var instance CORS
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

	return true, nil
}
