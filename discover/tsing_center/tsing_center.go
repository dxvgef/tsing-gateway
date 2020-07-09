package tsing_center

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"local/global"
)

// tsing center
type TsingCenter struct {
	Addr    string `json:"addr"`              // tsing center的地址
	Timeout int    `json:"timeout,omitempty"` // 连接超时时间(秒)
	Secret  string `json:"secret"`            // tsing center的连接密钥
}

// 新建探测器实例
func New(config string) (*TsingCenter, error) {
	var e TsingCenter
	e.Timeout = 5 // 默认超时时间
	err := json.Unmarshal([]byte(config), &e)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return &e, nil
}

// 获节点
func (self *TsingCenter) Fetch(serviceID string) (*url.URL, error) {
	var (
		node     global.NodeType
		endpoint strings.Builder
		req      *http.Request
		resp     *http.Response
		target   *url.URL
	)
	endpoint.WriteString(self.Addr)
	endpoint.WriteString("/services/")
	endpoint.WriteString(global.EncodeKey(serviceID))
	endpoint.WriteString("/select")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(self.Timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", endpoint.String(), nil)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	req.Header.Set("SECRET", self.Secret)
	var c http.Client
	resp, err = c.Do(req)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Warn().Err(err).Caller().Send()
		}
	}()

	if resp.Status != "200 OK" {
		err = errors.New("探测器没有获得正确的响应")
		log.Err(err).Str("response status", resp.Status).Caller().Send()
		return nil, err
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	err = node.UnmarshalJSON(body)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	endpoint.Reset()
	endpoint.WriteString("http://")
	endpoint.WriteString(node.IP)
	endpoint.WriteString(":")
	endpoint.WriteString(strconv.Itoa(int(node.Port)))
	target, err = url.Parse(endpoint.String())
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	return target, nil
}
