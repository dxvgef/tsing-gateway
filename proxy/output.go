package proxy

import (
	"encoding/json"

	"github.com/dxvgef/tsing-gateway/global"
)

type Data struct {
	Hosts     map[string]global.HostType     `json:"hosts,omitempty"`
	Routes    map[string]string              `json:"routes,omitempty"`
	Upstreams map[string]global.UpstreamType `json:"upstreams,omitempty"`
}

// 所有数据输出成json
func OutputJSON() ([]byte, error) {
	var data Data
	data.Hosts = make(map[string]global.HostType, global.SyncMapLen(&global.Hosts))
	data.Routes = make(map[string]string, global.SyncMapLen(&global.Routes))
	data.Upstreams = make(map[string]global.UpstreamType, global.SyncMapLen(&global.Upstreams))
	global.Hosts.Range(func(key, value interface{}) bool {
		data.Hosts[key.(string)] = value.(global.HostType)
		return true
	})
	global.Routes.Range(func(key, value interface{}) bool {
		data.Routes[key.(string)] = value.(string)
		return true
	})
	global.Upstreams.Range(func(key, value interface{}) bool {
		data.Upstreams[key.(string)] = value.(global.UpstreamType)
		return true
	})
	return json.Marshal(&data)
}
