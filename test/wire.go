// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"wuyuan_exam/internal/biz"
)

func initTaskUsercase(nodeId int64) (*biz.TaskUsercase, error) {
	panic(wire.Build(biz.ProviderSet, ProviderSet))
}
