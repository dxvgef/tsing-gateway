package wrr

import "errors"

var Pool *PoolType

// Weighted Round-Robin 加权轮循(与LVS类似)
type NodeType struct {
	Addr   string // 地址
	Weight int    // 权重值
}

type PoolType struct {
	nodes         []*NodeType // 节点列表
	nodeTotal     int         // 节点总数
	weightGCD     int         // 权总值最大公约数
	maxWeight     int         // 最大权重值
	lastIndex     int         // 最后命中的节点索引
	currentWeight int         // 当前权重值
}

func New() *PoolType {
	if Pool != nil {
		return Pool
	}
	var pool PoolType
	pool.nodes = []*NodeType{}
	pool.nodeTotal = 0
	pool.weightGCD = 0
	pool.maxWeight = 0
	pool.lastIndex = -1
	pool.currentWeight = 0
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
	p.calcMaxWeight(weight)
	p.nodes = append(p.nodes, node)
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
			p.calcMaxWeight(weight)
			return
		}
	}

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.calcMaxWeight(weight)
	p.nodes = append(p.nodes, node)
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
			p.calcAllGCD(p.nodeTotal)
		}
	}
}

// 选举出下一个命中的节点
func (p *PoolType) Next() string {
	if p.nodeTotal == 0 {
		return ""
	}

	if p.nodeTotal == 1 {
		return p.nodes[0].Addr
	}

	for {
		p.lastIndex = (p.lastIndex + 1) % p.nodeTotal
		if p.lastIndex == 0 {
			p.currentWeight = p.currentWeight - p.weightGCD
			if p.currentWeight <= 0 {
				p.currentWeight = p.maxWeight
				if p.currentWeight == 0 {
					return ""
				}
			}
		}

		if p.nodes[p.lastIndex].Weight >= p.currentWeight {
			return p.nodes[p.lastIndex].Addr
		}
	}
}

// 计算最大权重值
func (p *PoolType) calcMaxWeight(weight int) {
	if weight == 0 {
		return
	}
	if p.weightGCD == 0 {
		p.weightGCD = weight
		p.maxWeight = weight
		p.lastIndex = -1
		p.currentWeight = 0
		return
	}
	p.weightGCD = calcGCD(p.weightGCD, weight)
	if p.maxWeight < weight {
		p.maxWeight = weight
	}
}

// 计算所有节点权重值的最大公约数
func (p *PoolType) calcAllGCD(i int) int {
	if i == 1 {
		return p.nodes[0].Weight
	}
	return calcGCD(p.nodes[i-1].Weight, p.calcAllGCD(i-1))
}

// 计算两个权重值的最大公约数
func calcGCD(a, b int) int {
	if a < b {
		a, b = b, a // 交换a和b
	}
	if b == 0 {
		return a
	}
	return calcGCD(b, a%b)
}
