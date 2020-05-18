package api

// type Route struct{}
//
// func (self *Route) Put(ctx *tsing.Context) error {
// 	resp := make(map[string]string)
// 	var req struct {
// 		groupID    string
// 		path       string
// 		method     string
// 		upstreamID string
// 	}
// 	err := filter.MSet(
// 		filter.El(&req.groupID, filter.FromString(ctx.Post("group_id")).Required()),
// 		filter.El(&req.path, filter.FromString(ctx.Post("path")).Required()),
// 		filter.El(&req.method, filter.FromString(ctx.Post("method")).Required().EnumString()),
// 	)
// 	if err != nil {
// 		resp["error"] = err.Error()
// 		return JSON(ctx, 400, &resp)
// 	}
// 	if err = sa.PutUpstream(upstream.ID, upstreamStr); err != nil {
// 		resp["error"] = err.Error()
// 		return JSON(ctx, 500, &resp)
// 	}
// 	return Status(ctx, 204)
// }
//
// func (self *Upstream) Del(ctx *tsing.Context) error {
// 	resp := make(map[string]string)
// 	id := ctx.PathParams.Value("id")
// 	if id == "" {
// 		resp["error"] = "id参数不能为空"
// 		return JSON(ctx, 400, &resp)
// 	}
// 	err := sa.DelUpstream(id)
// 	if err != nil {
// 		resp["error"] = err.Error()
// 		return JSON(ctx, 500, &resp)
// 	}
// 	return Status(ctx, 204)
// }
