package etcd

// 获得锁
func (self *Etcd) getLock(target string) error {
	// var key strings.Builder
	// key.WriteString(self.KeyPrefix)
	// switch target {
	// case "host":
	// 	key.WriteString("/lock/hosts/")
	// case "upstream":
	// 	key.WriteString("/lock/upstreams/")
	// case "route":
	// 	key.WriteString("/lock/routes/")
	// default:
	// 	return errors.New("target值无效")
	// }
	//
	// ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer ctxCancel()
	// if _, err := self.client.Get(ctx, key.String()); err != nil {
	// 	log.Error().Caller().Msg(err.Error())
	// 	return err
	// }

	return nil
}
