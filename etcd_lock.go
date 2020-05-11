package main

// 获得锁
// func (p *Proxy) getLock(target string) (err error) {
// 	switch target {
// 	case "host":
// 		target = "/host_lock"
// 	case "upstream":
// 		target = "/upstream_lock"
// 	case "route":
// 		target = "/route_lock"
// 	}
// 	var key strings.Builder
// 	key.WriteString(global.Config.Etcd.KeyPrefix)
// 	key.WriteString(target)
// 	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer ctxCancel()
// 	if _, err = global.EtcdCli.Get(ctx, key.String()); err != nil {
// 		log.Error().Caller().Msg(err.Error())
// 		return err
// 	}
// 	key.Reset()
//
// 	// 写入路由
// 	for hostname, upstreamID := range p.Hosts {
// 		key.WriteString(global.Config.Etcd.KeyPrefix)
// 		key.WriteString("/hosts/")
// 		key.WriteString(hostname)
// 		ctx2, ctx2Cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		_, err = global.EtcdCli.Put(ctx2, key.String(), upstreamID)
// 		if err != nil {
// 			log.Error().Caller().Msg(err.Error())
// 			ctx2Cancel()
// 			return
// 		}
// 		key.Reset()
// 		ctx2Cancel()
// 	}
//
// 	return
// }
