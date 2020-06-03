package etcd

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadAllHosts(); err != nil {
		return
	}
	if err = self.LoadAllUpstreams(); err != nil {
		return
	}
	if err = self.LoadAllRoutes(); err != nil {
		return
	}
	return
}

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	if err = self.SaveAllHosts(); err != nil {
		return
	}
	if err = self.SaveAllUpstreams(); err != nil {
		return
	}
	if err = self.SaveAllRoutes(); err != nil {
		return
	}
	return
}
