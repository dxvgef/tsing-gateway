package wrr

import "errors"

var Pool *PoolType

// Weighted Round-Robin 加权轮循(与LVS类似)
type NodeType struct {
	Addr   string // 地址
	Weight int    // 权重值
}

type PoolType struct {
	nodes         map[string][]*NodeType // 节点列表
	nodeTotal     map[string]int         // 节点总数
	weightGCD     map[string]int         // 权总值最大公约数
	maxWeight     map[string]int         // 最大权重值
	lastIndex     map[string]int         // 最后命中的节点索引，初始值是-1
	currentWeight map[string]int         // 当前权重值
}

func Init() *PoolType {
	if Pool != nil {
		return Pool
	}
	var pool PoolType
	pool.nodes = map[string][]*NodeType{}
	pool.nodeTotal = map[string]int{}
	pool.weightGCD = map[string]int{}
	pool.maxWeight = map[string]int{}
	pool.lastIndex = map[string]int{}
	pool.currentWeight = map[string]int{}
	return &pool
}

func (p *PoolType) Add(upstreamID, addr string, weight int) error {
	if _, ok := p.nodes[upstreamID]; !ok {
		p.nodes[upstreamID] = []*NodeType{}
		p.lastIndex[upstreamID] = -1
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
	p.calcMaxWeight(upstreamID, weight)
	p.nodes[upstreamID] = append(p.nodes[upstreamID], node)
	p.nodeTotal[upstreamID]++
	return nil
}

func (p *PoolType) Set(upstreamID, addr string, weight int) {
	if p.nodes[upstreamID] == nil {
		p.nodes[upstreamID] = []*NodeType{}
		p.lastIndex[upstreamID] = -1
	}
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			p.nodes[upstreamID][i].Weight = weight
			p.calcMaxWeight(upstreamID, weight)
			return
		}
	}

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	p.calcMaxWeight(upstreamID, weight)
	p.nodes[upstreamID] = append(p.nodes[upstreamID], node)
	p.nodeTotal[upstreamID]++
}

// 节点总数
func (p *PoolType) Total(upstreamID string) int {
	return p.nodeTotal[upstreamID]
}

// 移除节点
func (p *PoolType) Remove(upstreamID, addr string) {
	for i := range p.nodes[upstreamID] {
		if p.nodes[upstreamID][i].Addr == addr {
			p.nodes[upstreamID] = append(p.nodes[upstreamID][:i], p.nodes[upstreamID][i+1:]...)
			p.nodeTotal[upstreamID]--
			p.calcAllGCD(upstreamID, p.nodeTotal[upstreamID])
		}
	}
}

// 选举出下一个命中的节点
func (p *PoolType) Next(upstreamID string) string {
	if p.nodeTotal[upstreamID] == 0 {
		return ""
	}

	if p.nodeTotal[upstreamID] == 1 {
		return p.nodes[upstreamID][0].Addr
	}

	for {
		p.lastIndex[upstreamID] = (p.lastIndex[upstreamID] + 1) % p.nodeTotal[upstreamID]
		if p.lastIndex[upstreamID] == 0 {
			p.currentWeight[upstreamID] = p.currentWeight[upstreamID] - p.weightGCD[upstreamID]
			if p.currentWeight[upstreamID] <= 0 {
				p.currentWeight[upstreamID] = p.maxWeight[upstreamID]
				if p.currentWeight[upstreamID] == 0 {
					return ""
				}
			}
		}

		if p.nodes[upstreamID][p.lastIndex[upstreamID]].Weight >= p.currentWeight[upstreamID] {
			return p.nodes[upstreamID][p.lastIndex[upstreamID]].Addr
		}
	}
}

// 计算最大权重值
func (p *PoolType) calcMaxWeight(upstreamID string, weight int) {
	if weight == 0 {
		return
	}
	if p.weightGCD[upstreamID] == 0 {
		p.weightGCD[upstreamID] = weight
		p.maxWeight[upstreamID] = weight
		p.lastIndex[upstreamID] = -1
		p.currentWeight[upstreamID] = 0
		return
	}
	p.weightGCD[upstreamID] = calcGCD(p.weightGCD[upstreamID], weight)
	if p.maxWeight[upstreamID] < weight {
		p.maxWeight[upstreamID] = weight
	}
}

// 计算所有节点权重值的最大公约数
func (p *PoolType) calcAllGCD(upstreamID string, i int) int {
	if i == 1 {
		return p.nodes[upstreamID][0].Weight
	}
	return calcGCD(p.nodes[upstreamID][i-1].Weight, p.calcAllGCD(upstreamID, i-1))
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
