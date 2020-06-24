package url_rewrite

import (
	"local/global"
	"net/http"

	"github.com/rs/zerolog/log"
)

// check_jwt
type CheckJWT struct {
	SourceChannel string `json:"source_channel"`           // token来源渠道，支持header、get、post、cookie
	SourceName    string `json:"source_name"`              // token来源参数名
	Alg           string `json:"alg"`                      // 用于本地校验的签名算法
	PrivateKey    string `json:"private_key,omitempty"`    // 本地校验用的密钥字符串(Base64)，用于本地校验
	RemoteURL     string `json:"remote_url,omitempty"`     // 发送给远程校验的URL
	RemoteChannel string `json:"remote_channel,omitempty"` // 发送给远程校验的通道
	RemoteName    string `json:"remote_name,omitempty"`    // 发送给远程校验的参数名
}

// 新建中间件实例
func New(config string) (*CheckJWT, error) {
	var instance CheckJWT
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}

func (self *CheckJWT) GetName() string {
	return "check_jwt"
}

// 中间件行为
func (self *CheckJWT) Action(req *http.Request) (bool, error) {
	log.Debug().Caller().Msg("check_jwt")
	return true, nil
}
