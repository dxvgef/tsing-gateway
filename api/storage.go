package api

import "github.com/dxvgef/tsing"

type Data struct{}

func (*Data) LoadAll(ctx *tsing.Context) error {
	if err := loadAll(); err != nil {
		return Status(ctx, 500)
	}
	return Status(ctx, 204)
}
func (*Data) SaveAll(ctx *tsing.Context) error {
	if err := saveAll(); err != nil {
		return Status(ctx, 500)
	}
	return Status(ctx, 204)
}

// 加载所有数据
func loadAll() (err error) {
	return sa.LoadAll()
}

// 保存所有数据
func saveAll() (err error) {
	return sa.SaveAll()
}
