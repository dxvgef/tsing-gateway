package etcd

// 加载所有数据
func (self *Etcd) LoadAll() (err error) {
	// 以下调用的函数中已经做了日志记录，所以在这里不再记录
	if err = self.LoadAllHost(); err != nil {
		return
	}
	if err = self.LoadAllService(); err != nil {
		return
	}
	if err = self.LoadAllRoute(); err != nil {
		return
	}
	return
}

// 存储所有数据
func (self *Etcd) SaveAll() (err error) {
	// 以下调用的函数中已经做了日志记录，所以在这里不再记录
	if err = self.SaveAllHost(); err != nil {
		return
	}
	if err = self.SaveAllService(); err != nil {
		return
	}
	if err = self.SaveAllRoute(); err != nil {
		return
	}
	return
}
