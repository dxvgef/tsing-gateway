package main

import (
	"encoding/json"
	"testing"

	"github.com/rs/zerolog/log"

	"github.com/dxvgef/tsing-gateway/global"
)

func TestRoute(t *testing.T) {
	var err error
	if err = global.LoadConfigFile(); err != nil {
		log.Fatal().Msg(err.Error())
	}
	if err = setLogger(); err != nil {
		log.Fatal().Msg(err.Error())
	}

	p := NewProxy()

	var upstream Upstream
	upstream.ID = "testUpstream"
	upstream.Middleware = append(upstream.Middleware, Configurator{
		Name:   "favicon",
		Config: `{"status": 204}`,
	})
	upstream.Explorer.Name = "coredns_etcd"
	upstream.Explorer.Config = `{"host":"test.uam.local"}`
	// 设置上游
	err = p.NewUpstream(upstream, false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 设置路由组
	routeGroup, err := p.SetRouteGroup("testGroup", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 设置路由
	err = routeGroup.SetRoute("/user/login", "get", "testUpstream", false)
	if err != nil {
		t.Error(err.Error())
		return
	}
	// 设置主机
	err = p.SetHost("127.0.0.1", routeGroup.ID, false)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// 序列化成json
	bb, err := json.Marshal(&p)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(bb))
}

func TestParse(t *testing.T) {
	var err error
	p := NewProxy()
	configStr := `{"hosts":{"127.0.0.1":"testGroup"},"route_groups":{"testGroup":{"/user/login":{"GET":"testUpstream"}}},"upstreams":{"testUpstream":{"id":"testUpstream","middleware":[{"name":"favicon","config":"{\"status\": 204}"}],"explorer":{"name":"coredns_etcd","config":"{\"host\":\"test.uam.local\"}"}}}}`
	err = json.Unmarshal([]byte(configStr), &p)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(p)
	bb, err := json.Marshal(&p)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(bb))
}
