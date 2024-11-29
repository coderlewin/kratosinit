package utils

import (
	"github.com/bwmarrin/snowflake"
	"sync"
)

var snowflakeNode *snowflake.Node

var once sync.Once

func GenSnowflakeID() int64 {
	once.Do(func() {
		snowflakeNode, _ = snowflake.NewNode(1)
	})
	return snowflakeNode.Generate().Int64()
}
