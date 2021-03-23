// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"wuyuan_exam/internal/biz"
	"wuyuan_exam/internal/data"
	"wuyuan_exam/internal/service"
)

func initTaskService(sqlConn string, nodeId int64, engine *gin.Engine) (*service.TaskService, error) {
	panic(wire.Build(data.ProviderSet, biz.ProviderSet, service.ProviderSet))
}
