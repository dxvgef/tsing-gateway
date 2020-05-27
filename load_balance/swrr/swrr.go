package swrr

import (
	"errors"
)

var Pool *PoolType

// Smooth Weighted Round-Robin 平滑加权轮循(与Nginx类似)
type PoolType struct {
	nodes     []*NodeType // 节点列表 map[upstreamID][addr]
	nodeTotal int         // 节点总量 map[upstreamID]
}

// 节点结构
type NodeType struct {
	Addr            string
	Weight          int
	CurrentWeight   int
	EffectiveWeight int
}

func New() *PoolType {
	if Pool != nil {
		return Pool
	}
	var pool PoolType
	pool.nodes = []*NodeType{}
	Pool = &pool
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
func (p *PoolType) Add(addr string, weight int) error {
	if p.nodes == nil {
		p.nodes = []*NodeType{}
	}
	for i := range p.nodes {
		if p.nodes[i].Addr == addr {
			return errors.New("节点地址已存在")
		}
	}
	node := NodeType{
		EffectiveWeight: weight,
		Weight:          weight,
		Addr:            addr,
	}
	p.nodes = append(p.nodes, &node)
	p.nodeTotal++
	return nil
}

// 设置节点
func (p *PoolType) Put(addr string, weight int) {
	if p.nodes == nil {
		p.nodes = []*NodeType{}
	}
	for i := range p.nodes {
		if p.nodes[i].Addr == addr {
			p.nodes[i].Weight = weight
			return
		}
	}
	node := NodeType{
		EffectiveWeight: weight,
		Weight:          weight,
		Addr:            addr,
	}
	p.nodes = append(p.nodes, &node)
	p.nodeTotal++
}

// 节点总数
func (p *PoolType) Total() int {
	return p.nodeTotal
}

// 移除指定地址的节点
func (p *PoolType) Remove(addr string) {
	for i := range p.nodes {
		if p.nodes[i].Addr == addr {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			p.nodeTotal--
		}
	}
}

// 选举出一个命中的节点
func (p *PoolType) Next() string {
	if p.nodeTotal == 0 {
		return ""
	}
	if p.nodeTotal == 1 {
		return p.nodes[0].Addr
	}
	var (
		addr   string
		target *NodeType
	)
	totalWeight := 0
	for i := range p.nodes {
		if p.nodes[i] == nil {
			continue
		}
		p.nodes[i].CurrentWeight += p.nodes[i].EffectiveWeight
		totalWeight += p.nodes[i].EffectiveWeight
		if p.nodes[i].EffectiveWeight < p.nodes[i].Weight {
			p.nodes[i].EffectiveWeight++
		}
		if target == nil || p.nodes[i].CurrentWeight > target.CurrentWeight {
			target = p.nodes[i]
			addr = p.nodes[i].Addr
		}
	}
	if target == nil {
		return ""
	}
	target.CurrentWeight -= totalWeight
	return addr
}

// 重置所有节点的状态
func (p *PoolType) Reset() {
	for i := range p.nodes {
		p.nodes[i].EffectiveWeight = p.nodes[i].Weight
		p.nodes[i].CurrentWeight = 0
	}
}
