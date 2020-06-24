package jwt_proxy

import (
	"errors"
	"io/ioutil"
	"local/global"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// JWTProxy
type JWTProxy struct {
	SourceType          string `json:"source_type"`                     // 来源类型，支持header、query、form、cookie
	SourceName          string `json:"source_name"`                     // 来源参数名
	UpstreamURL         string `json:"upstream_url"`                    // 上游URL
	SendType            string `json:"send_type"`                       // 发送给上游校验的类型，支持header、query、cookie
	SendMethod          string `json:"send_method"`                     // 发送给上游校验的HTTP方法，支持GET、HEAD、OPTIONS
	SendName            string `json:"send_name"`                       // 发送给上游校验的参数名
	UpstreamSuccessBody string `json:"upstream_success_body,omitempty"` // 校验成功的上游响应数据，用于处理只返回200状态码的API，如果留空则不校验body
}

// 新建中间件实例
func New(config string) (*JWTProxy, error) {
	var instance JWTProxy
	err := instance.UnmarshalJSON(global.StrToBytes(config))
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	instance.SourceType = strings.ToLower(instance.SourceType)
	instance.SendType = strings.ToLower(instance.SendType)
	instance.SendMethod = strings.ToUpper(instance.SendMethod)
	if instance.SourceType != "header" && instance.SourceType != "query" && instance.SourceType != "form" && instance.SourceType != "cookie" {
		return nil, errors.New("source_type参数值无效")
	}
	if instance.SendMethod != "GET" && instance.SendMethod != "HEAD" && instance.SendMethod != "OPTIONS" && instance.SendMethod != "TRACE" {
		return nil, errors.New("send_method参数值无效")
	}
	if instance.SourceName == "" {
		return nil, errors.New("source_name参数值无效")
	}
	if instance.UpstreamURL == "" {
		return nil, errors.New("upstream_url参数值无效")
	}
	if instance.SendName == "" {
		return nil, errors.New("send_name参数值无效")
	}
	if instance.SendType != "header" && instance.SendType != "query" && instance.SendType != "form" && instance.SendType != "cookie" {
		return nil, errors.New("send_type参数值无效")
	}
	return &instance, nil
}

func (self *JWTProxy) GetName() string {
	return "check_jwt"
}

// 中间件行为
func (self *JWTProxy) Action(resp http.ResponseWriter, req *http.Request) (bool, error) {
	log.Debug().Caller().Msg("check_jwt")
	// 获得jwt字符串
	jwtStr, err := self.getJWT(req)
	if err != nil {
		// 来自客户端数据，不记录日志
		resp.WriteHeader(http.StatusInternalServerError)
		log.Err(err).Caller().Send()
		return false, err
	}
	if jwtStr == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		return false, nil
	}

	// 发送jwt字符串
	var status int
	status, err = self.verityJWT(jwtStr)
	if err != nil {
		// 来自客户端数据，不记录日志
		resp.WriteHeader(status)
		log.Err(err).Caller().Send()
		return false, err
	}
	return true, nil
}

// 从客户端请求中获取JWT
func (self *JWTProxy) getJWT(req *http.Request) (string, error) {
	switch self.SourceType {
	case "header":
		return req.Header.Get(self.SourceName), nil
	case "query":
		if _, exist := req.URL.Query()[self.SourceName]; exist {
			return req.URL.Query()[self.SourceName][0], nil
		}
	case "form":
		if strings.HasPrefix(req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := req.ParseMultipartForm(http.DefaultMaxHeaderBytes); err != nil {
				log.Err(err).Caller().Send()
				return "", err
			}
		} else {
			if err := req.ParseForm(); err != nil {
				log.Err(err).Caller().Send()
				return "", err
			}
		}
		return req.FormValue(self.SourceName), nil
	case "cookie":
		ck, err := req.Cookie(self.SourceName)
		if err != nil {
			log.Err(err).Caller().Send()
			return "", err
		}
		if ck == nil {
			return "", nil
		}
		return ck.Value, nil
	}

	return "", nil
}

// 将JWT发送到远程服务器校验，并返回校验结果
// int == 0 && error == nil 表示校验通过
func (self *JWTProxy) verityJWT(jwtStr string) (int, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	// 根据发送类型准备http.Request
	switch self.SendType {
	case "header":
		req, err = http.NewRequest(self.SendMethod, self.UpstreamURL, nil)
		if err != nil {
			log.Err(err).Caller().Send()
			return http.StatusInternalServerError, err
		}
		req.Header.Set(self.SendName, jwtStr)
	case "query":
		endpoint, err := global.AddQueryValues(self.UpstreamURL, map[string]string{
			self.SendName: jwtStr,
		})
		if err != nil {
			log.Err(err).Caller().Send()
			return http.StatusInternalServerError, err
		}
		req, err = http.NewRequest(self.SendMethod, endpoint, nil)
		if err != nil {
			log.Err(err).Caller().Send()
			return http.StatusInternalServerError, err
		}
	case "cookie":
		req, err = http.NewRequest(self.SendMethod, self.UpstreamURL, nil)
		if err != nil {
			log.Err(err).Caller().Send()
			return http.StatusInternalServerError, err
		}
		req.AddCookie(&http.Cookie{
			Name:  self.SendName,
			Value: jwtStr,
		})
	}

	// 发送http.Request
	var (
		client http.Client
		body   []byte
	)
	resp, err = client.Do(req)
	if err != nil {
		log.Err(err).Caller().Send()
		return 0, err
	}

	// 处理远程响应
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Caller().Send()
		return http.StatusInternalServerError, err
	}
	bodyStr := global.BytesToStr(body)

	// 优先把两种可能成功的先判断
	switch resp.StatusCode {
	case http.StatusOK:
		if self.UpstreamSuccessBody != "" && self.UpstreamSuccessBody != bodyStr {
			return http.StatusForbidden, nil
		} else {
			return 0, nil
		}
	case http.StatusNoContent:
		return 0, nil
	}

	// 余下的全是不成功的处理
	if bodyStr == "" {
		bodyStr = http.StatusText(resp.StatusCode)
	}
	return resp.StatusCode, errors.New(bodyStr)
}
