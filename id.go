package main

import (
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/rs/zerolog/log"
)

// ID节点实例
var idNode *snowflake.Node

// 获得ID
func getID() string {
	var err error
	snowflake.Epoch = time.Now().Unix()
	idNode, err = snowflake.NewNode(0)
	if err != nil {
		log.Error().Msg(err.Error())
		return ""
	}
	return idNode.Generate().String()
}
