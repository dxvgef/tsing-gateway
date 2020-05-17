package api

import "github.com/dxvgef/tsing"

type Host struct{}

func (self *Host) Put(ctx *tsing.Context) error {
	err := putHost(ctx.Post("name"), ctx.Post("upstream_id"))
	if err != nil {
		return Status(ctx, 500)
	}
	return Status(ctx, 204)
}
func putHost(name, upstreamID string) error {
	return sa.PutHost(name, upstreamID)
}

func (self *Host) Del(ctx *tsing.Context) error {
	err := delHost(ctx.Post("name"))
	if err != nil {
		return Status(ctx, 500)
	}
	return Status(ctx, 204)
}
func delHost(name string) error {
	return sa.DelHost(name)
}
