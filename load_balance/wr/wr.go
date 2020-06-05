package wr

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

var (
	Inst         *InstType
	globalNodes  sync.Map // 节点列表 map[upstreamID][]*NodeType
	globalTotal  sync.Map // 节点总数 map[upstreamID]int
	globalWeight sync.Map // 节点权重总和 map[upstreamID]int
	globalRand   sync.Map // map[upstreamID]*rand.Rand
)

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
				if err = self.updateGlobalWeights(upstreamID, weight); err != nil {
					return
				}
				return self.Reset(upstreamID)
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
	if err = self.updateGlobalWeights(upstreamID, weight); err != nil {
		return
	}
	if err = self.Reset(upstreamID); err != nil {
		return
	}
	// 递增节点总数
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

// 移除所有节点
func (self *InstType) Remove(upstreamID, addr string) (err error) {
	mapValue, exist := globalNodes.Load(upstreamID)
	if !exist {
		return nil
	}
	nodes, ok := mapValue.([]*NodeType)
	if !ok {
		err = errors.New("节点类型断言失败")
		log.Err(err).Caller().Msg("移除节点")
		return
	}
	for k := range nodes {
		if nodes[k].Addr != addr {
			continue
		}
		newNodes := append(nodes[:k], nodes[k+1:]...)
		globalNodes.Store(upstreamID, newNodes)
		if err = self.updateGlobalWeights(upstreamID, -nodes[k].Weight); err != nil {
			return
		}
		if err = self.Reset(upstreamID); err != nil {
			return
		}
		return self.updateTotal(upstreamID, -1)
	}
	return
}

// 重设所有节点当前的权重值
func (self *InstType) Reset(upstreamID string) (err error) {
	if _, exist := globalRand.Load(upstreamID); !exist {
		return
	}
	globalRand.Store(upstreamID, rand.New(rand.NewSource(time.Now().UnixNano())))
	return nil
}

// 选举出下一个命中的节点
func (self *InstType) Next(upstreamID string) string {
	nodes := self.getNodes(upstreamID)
	if nodes == nil {
		return ""
	}
	nodeTotal := len(nodes)
	if nodeTotal == 1 {
		return nodes[0].Addr
	}
	var rnd *rand.Rand
	mapValue, exist := globalRand.Load(upstreamID)
	if !exist {
		if r, ok := mapValue.(*rand.Rand); ok {
			rnd = r
		}
	}
	if rnd == nil {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
		globalRand.Store(upstreamID, rnd)
	}

	randomWeight := rnd.Intn(self.getGlobalWeight(upstreamID))
	for k := range nodes {
		randomWeight = randomWeight - nodes[k].Weight
		if randomWeight <= 0 {
			return nodes[k].Addr
		}
	}
	return ""
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

// 更新节点统计总数
func (self *InstType) updateGlobalWeights(upstreamID string, count int) (err error) {
	mapValue, exist := globalWeight.Load(upstreamID)
	if !exist {
		globalWeight.Store(upstreamID, 0)
		return nil
	}

	total, ok := mapValue.(int)
	if !ok {
		err = errors.New("类型断言失败")
		log.Err(err).Caller().Msg("更新权重值总数")
		return err
	}
	total += count
	globalWeight.Store(upstreamID, total)
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

// 获得权重值总和
func (self *InstType) getGlobalWeight(upstreamID string) int {
	mapValue, exist := globalWeight.Load(upstreamID)
	if !exist {
		return 0
	}
	total, ok := mapValue.(int)
	if !ok {
		return 0
	}
	return total
}
