package swrr

import (
	"errors"
)

var Pool *PoolType

// Smooth Weighted Round-Robin 平滑加权轮循(与Nginx类似)
type PoolType struct {
	nodes     map[string][]*NodeType // 节点列表 map[upstreamID][addr]
	nodeTotal map[string]int         // 节点总量 map[upstreamID]
}

// 节点结构
type NodeType struct {
	Addr            string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func Init() *PoolType {
	if Pool != nil {
		return Pool
	}
	Pool.nodes = map[string][]*NodeType{}
	Pool.nodeTotal = map[string]int{}
	return Pool
}

// // 降权
// func (n *NodeType) Reduce() {
// 	n.EffectiveWeight -= n.Weight
// 	if n.EffectiveWeight < 0 {
// 		n.EffectiveWeight = 0
// 	}
// }

// 添加节点
func (p *PoolType) Add(upstreamID, addr string, weight int) error {
	if _, ok := p.nodes[upstreamID]; !ok {
		p.nodes[upstreamID] = []*NodeType{}
	}
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			return errors.New("节点地址已存在")
		}
	}
	node := NodeType{
		EffectiveWeight: weight,
		Weight:          weight,
		Addr:            addr,
	}
	p.nodes[upstreamID] = append(p.nodes[upstreamID], &node)
	p.nodeTotal[upstreamID]++
	return nil
}

// 设置节点
func (p *PoolType) Put(upstreamID, addr string, weight int) {
	if _, ok := p.nodes[upstreamID]; !ok {
		p.nodes[upstreamID] = []*NodeType{}
	}
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			p.nodes[upstreamID][i].Weight = weight
			return
		}
	}
	node := NodeType{
		EffectiveWeight: weight,
		Weight:          weight,
		Addr:            addr,
	}
	p.nodes[upstreamID] = append(p.nodes[upstreamID], &node)
	p.nodeTotal[upstreamID]++
}

// 节点总数
func (p *PoolType) Total(upstreamID string) int {
	return p.nodeTotal[upstreamID]
}

// 移除指定地址的节点
func (p *PoolType) Remove(upstreamID, addr string) {
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			p.nodes[upstreamID] = append(p.nodes[upstreamID][:i], p.nodes[upstreamID][i+1:]...)
			p.nodeTotal[upstreamID]--
		}
	}
}

// 选举出一个命中的节点
func (p *PoolType) Next(upstreamID string) string {
	if p.nodeTotal[upstreamID] == 0 {
		return ""
	}
	if p.nodeTotal[upstreamID] == 1 {
		return p.nodes[upstreamID][0].Addr
	}
	var (
		addr   string
		target *NodeType
	)
	totalWeight := 0
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i] == nil {
			continue
		}
		p.nodes[upstreamID][i].CurrentWeight += p.nodes[upstreamID][i].EffectiveWeight
		totalWeight += p.nodes[upstreamID][i].EffectiveWeight
		if p.nodes[upstreamID][i].EffectiveWeight < p.nodes[upstreamID][i].Weight {
			p.nodes[upstreamID][i].EffectiveWeight++
		}
		if target == nil || p.nodes[upstreamID][i].CurrentWeight > target.CurrentWeight {
			target = p.nodes[upstreamID][i]
			addr = p.nodes[upstreamID][i].Addr
		}
	}
	if target == nil {
		return ""
	}
	target.CurrentWeight -= totalWeight
	return addr
}

// 重置所有节点的状态
func (p *PoolType) Reset(upstreamID string) {
	for i := range p.nodes[upstreamID] {
		p.nodes[upstreamID][i].EffectiveWeight = p.nodes[upstreamID][i].Weight
		p.nodes[upstreamID][i].CurrentWeight = 0
	}
}
