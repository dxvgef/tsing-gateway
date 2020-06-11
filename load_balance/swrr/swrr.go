package swrr

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

var Inst *InstType
var globalNodes sync.Map // 节点列表 key=upstreamID, value=[]*NodeType
var globalTotal sync.Map // 节点总量 key=upstreamID, value=int

type InstType struct{}

// 节点结构
type NodeType struct {
	Addr            string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func Init() *InstType {
	if Inst != nil {
		return Inst
	}
	return &InstType{}
}

// // 降权
// func (n *NodeType) Reduce() {
// 	n.EffectiveWeight -= n.Weight
// 	if n.EffectiveWeight < 0 {
// 		n.EffectiveWeight = 0
// 	}
// }

// 设置节点
func (self *InstType) Set(upstreamID, addr string, weight int) (err error) {
	var nodes []*NodeType
	mapValue, exist := globalNodes.Load(upstreamID)
	if exist {
		var ok bool
		if nodes, ok = mapValue.([]*NodeType); !ok {
			err = errors.New("类型断言失败")
			log.Err(err).Caller().Msg("设置节点")
			return
		}
	}
	for k := range nodes {
		// 如果节点已存在，则直接更新
		if nodes[k].Addr == addr {
			if nodes[k].Weight != weight {
				nodes[k].Weight = weight
				globalNodes.Store(upstreamID, nodes)
				return self.reset(upstreamID)
			}
			return
		}
	}

	// 插入节点
	nodes = append(nodes, &NodeType{
		Addr:   addr,
		Weight: weight,
	})
	globalNodes.Store(upstreamID, nodes)
	if err = self.reset(upstreamID); err != nil {
		return
	}
	// 递增节点总数
	return self.updateTotal(upstreamID, 1)
}

// 移除指定地址的节点
func (self *InstType) Remove(upstreamID, addr string) (err error) {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	if addr == "" {

		globalNodes.Delete(upstreamID)
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("移除节点")
		return
	}
	for k := range nodes {
		if nodes[k].Addr != addr {
			continue
		}
		newNodes := append(nodes[:k], nodes[k+1:]...)
		globalNodes.Store(upstreamID, newNodes)
		if err = self.reset(upstreamID); err != nil {
			return
		}
		return self.updateTotal(upstreamID, -1)
	}
	return
}

// 节点总数
func (self *InstType) Total(upstreamID string) int {
	mapValue, exist := globalTotal.Load(upstreamID)
	if !exist {
		return 0
	}
	total, ok := mapValue.(int)
	if !ok {
		return 0
	}
	return total
}

// 选举出一个命中的节点
func (self *InstType) Next(upstreamID string) string {
	nodes := self.getNodes(upstreamID)
	if nodes == nil {
		return ""
	}
	nodeTotal := len(nodes)
	if nodeTotal == 1 {
		return nodes[0].Addr
	}
	var (
		addr   string
		target *NodeType
	)
	totalWeight := 0
	for i := range nodes {
		nodes[i].CurrentWeight += nodes[i].EffectiveWeight
		totalWeight += nodes[i].EffectiveWeight
		if nodes[i].EffectiveWeight < nodes[i].Weight {
			nodes[i].EffectiveWeight++
		}
		if target == nil || nodes[i].CurrentWeight > target.CurrentWeight {
			target = nodes[i]
			addr = nodes[i].Addr
		}
	}

	globalNodes.Store(upstreamID, nodes)

	if target == nil {
		return ""
	}
	target.CurrentWeight -= totalWeight
	return addr
}

// 重置所有节点的状态
func (self *InstType) reset(upstreamID string) (err error) {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("重置节点")
		return
	}
	for k := range nodes {
		nodes[k].EffectiveWeight = nodes[k].Weight
		nodes[k].CurrentWeight = 0
	}
	return nil
}

// 更新节点统计总数
func (self *InstType) updateTotal(upstreamID string, count int) (err error) {
	mapValue, exist := globalTotal.Load(upstreamID)
	if !exist {
		globalTotal.Store(upstreamID, 0)
		return nil
	}

	total, ok := mapValue.(int)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("更新节点总数")
		return err
	}
	total += count
	globalTotal.Store(upstreamID, total)
	return nil
}

// 判断节点是否存在
func (self *InstType) nodeExist(upstreamID, addr string) (exist bool) {
	if _, exist = globalNodes.Load(upstreamID); !exist {
		return
	}
	globalNodes.Range(func(key, value interface{}) bool {
		if key.(string) == upstreamID {
			nodes, ok := value.([]*NodeType)
			if !ok {
				return true
			}
			for k := range nodes {
				if nodes[k].Addr == addr {
					exist = true
					return false
				}
			}
		}
		return true
	})
	return
}

// 获得节点
func (self *InstType) getNodes(upstreamID string) []*NodeType {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		return nil
	}
	return nodes
}
