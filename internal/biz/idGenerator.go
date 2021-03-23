package biz

import (
	"github.com/bwmarrin/snowflake"
)

type IdGenerator struct {
	nodeId int64
}

func NewIdGenerator(nodeId int64) *IdGenerator {
	return &IdGenerator{nodeId: nodeId}
}

func (this *IdGenerator) Generate() (int64, error) {
	node, err := snowflake.NewNode(this.nodeId)
	if err != nil {
		return -1, err
	}

	return node.Generate().Int64(), nil
}
