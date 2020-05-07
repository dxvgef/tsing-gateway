package main

import (
	"encoding/json"
	"testing"
)

func BenchmarkEcho(b *testing.B) {
	if err := setLogger(); err != nil {
		b.Fatal(err.Error())
		return
	}
	var endpoints []Endpoint
	var err error
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10080",
		Weight: 100,
	})
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10082",
		Weight: 100,
	})
	endpoints = append(endpoints, Endpoint{
		Addr:   "127.0.0.1:10084",
		Weight: 100,
	})

	proxy := New()
	// 添加上游及端点
	if err = proxy.setUpstream(Upstream{
		ID:        "userLogin",
		Endpoints: endpoints,
	}); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "userRegister",
		Endpoints: endpoints,
	}); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "user",
		Endpoints: endpoints,
	}); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = proxy.setUpstream(Upstream{
		ID:        "root",
		Endpoints: endpoints,
	}); err != nil {
		b.Fatal(err.Error())
		return
	}
	// 添加路由组
	routeGroup, err := proxy.newRouteGroup("uam_v1_routes")
	if err != nil {
		b.Fatal(err.Error())
		return
	}
	// 在路由组内写入路由规则
	if err = routeGroup.setRoute("/user/login", "GET", "userLogin"); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = routeGroup.setRoute("/user/register", "GET", "userRegister"); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = routeGroup.setRoute("/user/*", "GET", "user"); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = routeGroup.setRoute("/", "GET", "root"); err != nil {
		b.Fatal(err.Error())
		return
	}
	// 添加主机
	if err = proxy.setHost("127.0.0.1", "uam_v1_routes"); err != nil {
		b.Fatal(err.Error())
		return
	}
	if err = proxy.setHost("192.168.50.144", "uam_v1_routes"); err != nil {
		b.Fatal(err.Error())
		return
	}

	b.Log(proxy.hosts)
	hostsJSON, err := json.Marshal(&proxy.hosts)
	if err != nil {
		b.Fatal(err.Error())
	}
	b.Log("hosts", string(hostsJSON))

	b.Log(proxy.upstreams)
	upstreamsJSON, err := json.Marshal(&proxy.upstreams)
	if err != nil {
		b.Fatal(err.Error())
	}
	b.Log("upstreams", string(upstreamsJSON))

	b.Log(proxy.routeGroups)
	routeGroupsJSON, err := json.Marshal(&proxy.routeGroups)
	if err != nil {
		b.Fatal(err.Error())
	}
	b.Log("routeGroups", string(routeGroupsJSON))

	// req, err := http.NewRequest("GET", "http://127.0.0.1:10080/", nil)
	// if err != nil {
	// 	b.Fatal(err.Error())
	// 	return
	// }
	// proxy.ServeHTTP(httptest.NewRecorder(), req)
}
