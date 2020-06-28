package url_rewrite

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"local/global"
)

// url_rewrite
type URLRewrite struct {
	Prefix  map[string]string `json:"prefix,omitempty"`
	Suffix  map[string]string `json:"suffix,omitempty"`
	Replace map[string]string `json:"replace,omitempty"`
}

// 新建中间件实例
func New(config string) (*URLRewrite, error) {
	var instance URLRewrite
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &instance, nil
}

func (self *URLRewrite) GetName() string {
	return "url_rewrite"
}

// 中间件行为
func (self *URLRewrite) Action(_ http.ResponseWriter, req *http.Request) (bool, error) {
	// 前缀重写
	for k := range self.Prefix {
		if strings.HasPrefix(req.URL.Path, k) {
			req.URL.Path = strings.Replace(req.URL.Path, k, self.Prefix[k], 1)
			req.RequestURI = req.URL.RequestURI()
		}
	}
	// 后缀重写
	for k := range self.Suffix {
		if strings.HasSuffix(req.URL.Path, k) {
			req.URL.Path = strings.TrimSuffix(req.URL.Path, k)
			req.URL.Path += self.Suffix[k]
			req.RequestURI = req.URL.RequestURI()
		}
	}
	// 替换重写
	for k := range self.Replace {
		req.URL.Path = strings.ReplaceAll(req.URL.Path, k, self.Replace[k])
		req.RequestURI = req.URL.RequestURI()
	}
	return false, nil
}
