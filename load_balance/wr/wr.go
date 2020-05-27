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
	nodes        []*NodeType // 节点列表
	nodeTotal    int         // 节点总数
	sumOfWeights int         // 所有节点权重值的总和
	rnd          *rand.Rand
}

func New() *PoolType {
	if Pool != nil {
		return Pool
	}
	var pool PoolType
	pool.nodes = []*NodeType{}
	pool.nodeTotal = 0
	pool.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	Pool = &pool
	return Pool
}

func (p *PoolType) Add(addr string, weight int) error {
	for i := range p.nodes {
		if p.nodes[i].Addr == addr {
			return errors.New("节点地址已存在")
		}
	}

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.nodes = append(p.nodes, node)
	p.sumOfWeights += weight
	p.nodeTotal++
	return nil
}

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

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.nodes = append(p.nodes, node)
	p.sumOfWeights += weight
	p.nodeTotal++
}

// 节点总数
func (p *PoolType) Total() int {
	return p.nodeTotal
}

// 移除所有节点
func (p *PoolType) Remove(addr string) {
	for i := range p.nodes {
		if p.nodes[i].Addr == addr {
			p.nodes = append(p.nodes[:i], p.nodes[i+1:]...)
			p.nodeTotal--
			p.sumOfWeights -= p.nodes[i].Weight
		}
	}
}

// 重设所有节点当前的权重值
func (p *PoolType) Reset() {
	p.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// 选举出下一个命中的节点
func (p *PoolType) Next() string {
	if p.nodeTotal == 0 {
		return ""
	}

	if p.nodeTotal == 1 {
		return p.nodes[0].Addr
	}

	randomWeight := p.rnd.Intn(p.sumOfWeights)
	for k := range p.nodes {
		randomWeight = randomWeight - p.nodes[k].Weight
		if randomWeight <= 0 {
			return p.nodes[k].Addr
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
