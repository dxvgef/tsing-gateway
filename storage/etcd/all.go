package etcd

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	if err = self.LoadAllHost(); err != nil {
		return
	}
	if err = self.LoadAllUpstream(); err != nil {
		return
	}
	if err = self.LoadAllRoute(); err != nil {
		return
	}
	return
}

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	if err = self.SaveAllHost(); err != nil {
		return
	}
	if err = self.SaveAllUpstream(); err != nil {
		return
	}
	if err = self.SaveAllRoute(); err != nil {
		return
	}
	return
}
