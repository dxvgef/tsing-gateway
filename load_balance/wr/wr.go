package wr

import (
	"errors"
	"math/rand"
	"time"
)

var Pool *PoolType

// Weighted Random 加权随机(与LVS类似)
type NodeType struct {
	Addr   string // 地址
	Weight int    // 权重值
}

type PoolType struct {
	nodes        map[string][]*NodeType // 节点列表 map[upstreamID]
	nodeTotal    map[string]int         // 节点总数 map[upstreamID]
	sumOfWeights map[string]int         // 所有节点权重值的总和 map[upstreamID]
	rnd          map[string]*rand.Rand  // map[upstreamID]
}

func Init() *PoolType {
	if Pool != nil {
		return Pool
	}
	var pool PoolType
	pool.nodes = map[string][]*NodeType{}
	pool.nodeTotal = map[string]int{}
	pool.sumOfWeights = map[string]int{}
	pool.rnd = map[string]*rand.Rand{}
	return &pool
}

func (p *PoolType) Add(upstreamID, addr string, weight int) error {
	if _, ok := p.nodes[upstreamID]; !ok {
		p.nodes[upstreamID] = []*NodeType{}
	}
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			return errors.New("节点地址已存在")
		}
	}

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.nodes[upstreamID] = append(p.nodes[upstreamID], node)
	p.sumOfWeights[upstreamID] += weight
	p.nodeTotal[upstreamID]++
	return nil
}

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
	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.nodes[upstreamID] = append(p.nodes[upstreamID], node)
	p.sumOfWeights[upstreamID] += weight
	p.nodeTotal[upstreamID]++
}

// 节点总数
func (p *PoolType) Total(upstreamID string) int {
	return p.nodeTotal[upstreamID]
}

// 移除所有节点
func (p *PoolType) Remove(upstreamID, addr string) {
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			p.nodes[upstreamID] = append(p.nodes[upstreamID][:i], p.nodes[upstreamID][i+1:]...)
			p.nodeTotal[upstreamID]--
			p.sumOfWeights[upstreamID] -= p.nodes[upstreamID][i].Weight
		}
	}
}

// 重设所有节点当前的权重值
func (p *PoolType) Reset(upstreamID string) {
	p.rnd[upstreamID] = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 选举出下一个命中的节点
func (p *PoolType) Next(upstreamID string) string {
	if p.nodeTotal[upstreamID] == 0 {
		return ""
	}

	if p.nodeTotal[upstreamID] == 1 {
		return p.nodes[upstreamID][0].Addr
	}

	if p.rnd[upstreamID] == nil {
		p.rnd[upstreamID] = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	randomWeight := p.rnd[upstreamID].Intn(p.sumOfWeights[upstreamID])
	for k := range p.nodes[upstreamID] {
		randomWeight = randomWeight - p.nodes[upstreamID][k].Weight
		if randomWeight <= 0 {
			return p.nodes[upstreamID][k].Addr
		}
	}
	return ""
}

// 计算gcd的值
func gcd(x, y int) int {
	var t int
	for {
		t = x % y
		if t > 0 {
			x = y
			y = t
		} else {
			return y
		}
	}
}
