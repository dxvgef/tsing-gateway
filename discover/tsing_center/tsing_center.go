package tsing_center

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"

	tsingCenterCli "github.com/dxvgef/tsing-center-go"

	"github.com/rs/zerolog/log"
)

// tsing center
type TsingCenter struct {
	Addr      string `json:"addr"`              // tsing center的地址
	Timeout   uint   `json:"timeout,omitempty"` // 连接超时时间(秒)
	Secret    string `json:"secret"`            // tsing center的连接密钥
	ServiceID string `json:"service_id"`        // tsing center里的服务ID
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
func (self *TsingCenter) Fetch() (*url.URL, error) {
	var (
		err      error
		endpoint strings.Builder
		target   *url.URL
		status   int
		node     tsingCenterCli.Node

		cli *tsingCenterCli.Client
	)

	cli, err = tsingCenterCli.New(tsingCenterCli.Config{
		Addr:    self.Addr,
		Secret:  self.Secret,
		Timeout: self.Timeout,
	})
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

	node, status, err = cli.DiscoverService(self.ServiceID)
	if err != nil {
		err = errors.New("服务节点发现失败")
		log.Err(err).Str("serviceID", self.ServiceID).Str("discover", "tsing center").Int("status", status).Caller().Send()
		return nil, err
	}

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
