package global

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

// ID节点实例
var IDNode *snowflake.Node

// 获得ID string值
func GetIDStr() string {
	var err error
	snowflake.Epoch = time.Now().Unix()
	IDNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return ""
	}
	return IDNode.Generate().String()
}

// 获得ID int64值
func GetIDInt64() int64 {
	var err error
	snowflake.Epoch = time.Now().Unix()
	IDNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return 0
	}
	return IDNode.Generate().Int64()
}
