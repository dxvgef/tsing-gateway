package main

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

// ID节点实例
var idNode *snowflake.Node

// 获得ID string值
func getIDStr() string {
	var err error
	snowflake.Epoch = time.Now().Unix()
	idNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return ""
	}
	return idNode.Generate().String()
}

// 获得ID int64值
func getIDInt64() int64 {
	var err error
	snowflake.Epoch = time.Now().Unix()
	idNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return 0
	}
	return idNode.Generate().Int64()
}
