package tsing_center

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

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
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}
	req.Header.Set("SECRET", self.Secret)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Err(err).Caller().Send()
		return nil, err
	}

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
