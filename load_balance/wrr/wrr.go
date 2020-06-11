package wrr

import (
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	Inst          *InstType
	globalNodes   sync.Map // 节点列表
	globalTotal   sync.Map // 节点总数
	weightGCD     sync.Map // 权总值最大公约数
	maxWeight     sync.Map // 最大权重值
	lastIndex     sync.Map // 最后命中的节点索引，初始值是-1
	currentWeight sync.Map // 当前权重值
)

// Weighted Round-Robin 加权轮循(与LVS类似)
type NodeType struct {
	Addr   string // 地址
	Weight int    // 权重值
}

type InstType struct{}

func Init() *InstType {
	if Inst != nil {
		return Inst
	}
	return &InstType{}
}
func (self *InstType) Set(upstreamID, addr string, weight int) (err error) {
	var nodes []*NodeType
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		globalNodes.Store(upstreamID, []*NodeType{})
		lastIndex.Store(upstreamID, -1)
	} else {
		var ok bool
		if nodes, ok = mapValue.([]*NodeType); !ok {
			err = errors.New("类型断言失败")
			log.Err(err).Caller().Msg("设置节点")
			return
		}
	}
	for i := range nodes {
		if nodes[i].Addr == addr {
			nodes[i].Weight = weight
			self.calcMaxWeight(upstreamID, weight)
			return
		}
	}

	node := &NodeType{
		Addr:   addr,
		Weight: weight,
	}
	self.calcMaxWeight(upstreamID, weight)
	nodes = append(nodes, node)
	globalNodes.Store(upstreamID, nodes)
	return self.updateTotal(upstreamID, 1)
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

// 移除节点
func (self *InstType) Remove(upstreamID, addr string) (err error) {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("移除节点")
		return
	}
	mapValue, exist = globalTotal.Load(upstreamID)
	if !exist {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("移除节点")
		return
	}
	total, ok := mapValue.(int)
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
		self.calcAllGCD(upstreamID, total)
		return self.updateTotal(upstreamID, -1)
	}
	return
}

// 选举出下一个命中的节点
func (self *InstType) Next(upstreamID string) string {
	var (
		nodes    = self.getNodes(upstreamID)
		nodesLen = len(nodes)
	)
	switch nodesLen {
	case 0:
		return ""
	case 1:
		return nodes[0].Addr
	}

	var last, total, cw, gcd int
	for {
		last = self.getLastIndex(upstreamID)
		total = self.Total(upstreamID)
		last = (last + 1) % total
		lastIndex.Store(upstreamID, last)
		if last == 0 {
			cw = self.getCurrentWeight(upstreamID)
			gcd = self.getWeightGCD(upstreamID)
			newCW := cw - gcd
			currentWeight.Store(upstreamID, newCW)
			if newCW <= 0 {
				newCW = self.getMaxWeight(upstreamID)
				currentWeight.Store(upstreamID, newCW)
				if newCW == 0 {
					return ""
				}
			}
		}
		cw = self.getCurrentWeight(upstreamID)
		if nodes[last].Weight >= cw {
			return nodes[last].Addr
		}
	}
}

// 计算最大权重值
func (self *InstType) calcMaxWeight(upstreamID string, weight int) {
	if weight == 0 {
		return
	}
	mapValue, exist := weightGCD.Load(upstreamID)
	if !exist {
		weightGCD.Store(upstreamID, weight)
		maxWeight.Store(upstreamID, weight)
		lastIndex.Store(upstreamID, -1)
		currentWeight.Store(upstreamID, 0)
		return
	}
	cgd, ok := mapValue.(int)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("计算最大权重值")
		return
	}
	if cgd == 0 {
		weightGCD.Store(upstreamID, weight)
		maxWeight.Store(upstreamID, weight)
		lastIndex.Store(upstreamID, -1)
		currentWeight.Store(upstreamID, 0)
		return
	}
	weightGCD.Store(upstreamID, calcGCD(cgd, weight))
	mapValue, exist = maxWeight.Load(upstreamID)
	if !exist {
		return
	}
	max, ok := mapValue.(int)
	if !ok {
		return
	}
	if max < weight {
		maxWeight.Store(upstreamID, weight)
	}
}

// 计算所有节点权重值的最大公约数
func (self *InstType) calcAllGCD(upstreamID string, i int) int {
	nodes := self.getNodes(upstreamID)
	if i == 1 {
		return nodes[0].Weight
	}
	return calcGCD(nodes[i-1].Weight, self.calcAllGCD(upstreamID, i-1))
}

func (self *InstType) getLastIndex(upstreamID string) int {
	mapValue, exist := lastIndex.Load(upstreamID)
	if !exist {
		globalTotal.Store(upstreamID, -1)
		return -1
	}

	index, ok := mapValue.(int)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("获得节点列表")
		return -1
	}
	return index
}

// 获得节点
func (self *InstType) getNodes(upstreamID string) []*NodeType {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("获得节点列表")
		return nil
	}
	return nodes
}

// 更新节点统计总数
func (self *InstType) updateTotal(upstreamID string, count int) (err error) {
	mapValue, exist := globalTotal.Load(upstreamID)
	if !exist {
		globalTotal.Store(upstreamID, count)
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

func (self *InstType) getCurrentWeight(upstreamID string) int {
	mapValue, exist := currentWeight.Load(upstreamID)
	if !exist {
		return 0
	}

	weight, ok := mapValue.(int)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("获得当前权重值")
		return 0
	}
	return weight
}

func (self *InstType) getWeightGCD(upstreamID string) int {
	mapValue, exist := weightGCD.Load(upstreamID)
	if !exist {
		return 0
	}

	value, ok := mapValue.(int)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("获得权重GCD")
		return 0
	}
	return value
}

func (self *InstType) getMaxWeight(upstreamID string) int {
	mapValue, exist := maxWeight.Load(upstreamID)
	if !exist {
		return 0
	}

	value, ok := mapValue.(int)
	if !ok {
		log.Err(errors.New("类型断言失败")).Caller().Msg("获得最大权重值")
		return 0
	}
	return value
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
