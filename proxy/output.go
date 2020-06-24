package proxy

import (
	"encoding/json"

	"local/global"
)

type Data struct {
	Hosts    map[string]global.HostType    `json:"hosts,omitempty"`
	Routes   map[string]string             `json:"routes,omitempty"`
	Services map[string]global.ServiceType `json:"services,omitempty"`
}

// 所有数据输出成json
func OutputJSON() ([]byte, error) {
	var data Data
	data.Hosts = make(map[string]global.HostType, global.SyncMapLen(&global.Hosts))
	data.Routes = make(map[string]string, global.SyncMapLen(&global.Routes))
	data.Services = make(map[string]global.ServiceType, global.SyncMapLen(&global.Services))
	global.Hosts.Range(func(key, value interface{}) bool {
		data.Hosts[key.(string)] = value.(global.HostType)
		return true
	})
	global.Routes.Range(func(key, value interface{}) bool {
		data.Routes[key.(string)] = value.(string)
		return true
	})
	global.Services.Range(func(key, value interface{}) bool {
		data.Services[key.(string)] = value.(global.ServiceType)
		return true
	})
	return json.Marshal(&data)
}
