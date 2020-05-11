package main

import (
	"encoding/json"
	"flag"
	"testing"

	"github.com/dxvgef/tsing-gateway/global"
)

// func TestLock(t *testing.T) {
// 	go func() {
// 		writerID := global.GetIDInt64()
// 		for {
// 			t.Log("尝试写值")
// 			time.Sleep(2 * time.Second)
// 			t.Log("写值成功")
// 		}
// 	}()
// 	go func() {
// 		readerID := global.GetIDInt64()
// 		for {
// 			t.Log("尝试取值")
// 			time.Sleep(2 * time.Second)
// 			t.Log("取值成功")
// 		}
// 	}()
// 	select {}
// }

func TestRoute(t *testing.T) {
	var err error
	var configFile string
	flag.StringVar(&configFile, "c", "./config.local.yml", "配置文件路径")
	flag.Parse()
	if err = global.LoadConfigFile(configFile); err != nil {
		t.Error(err.Error())
		return
	}
	if err = setLogger(); err != nil {
		t.Error(err.Error())
		return
	}

	if err = global.SetEtcdCli(); err != nil {
		t.Error(err.Error())
		return
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

	// 设置上游
	upstream = Upstream{}
	upstream.ID = "test2Upstream"
	upstream.Explorer.Name = "coredns_etcd"
	upstream.Explorer.Config = `{"host":"test2.uam.local"}`
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

	// 将配置保存到etcd
	err = p.SaveDataToEtcd()
	if err != nil {
		t.Error(err.Error())
		return
	}

}
