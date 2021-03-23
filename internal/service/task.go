package service

import (
	"errors"
	ginserver "gin-server"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strconv"
	"wuyuan_exam/internal/biz"
)

type TaskService struct {
	engine       *gin.Engine
	taskUsercase *biz.TaskUsercase
}

func NewTaskService(engine *gin.Engine, usercase *biz.TaskUsercase) *TaskService {
	return &TaskService{
		engine:       engine,
		taskUsercase: usercase,
	}
}

func (taskService *TaskService) CreateTaskGroup() *gin.RouterGroup {
	handlerConfig := ginserver.Config{
		ErrorHandleFunc: func(context *gin.Context, err error) {
			context.AbortWithError(401, err)
		},
		TokenKey: "github.com/go-oauth2/gin-server/access-token",
		Skipper: func(_ *gin.Context) bool {
			return false
		},
	}

	taskGroup := taskService.engine.Group("/task")
	taskGroup.Use(ginserver.HandleTokenVerify(handlerConfig))
	taskGroup.POST("/", func(c *gin.Context) {
		var task biz.Task
		err := c.BindJSON(&task)
		if err != nil {
			WriteError(c, err)
			return
		}

		_, err = taskService.taskUsercase.Create(c, &task)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(200, gin.H{})
	})
	taskGroup.PUT("/:id/description", func(c *gin.Context) {
		var (
			description string
		)
		idStr, ok := c.Params.Get("id")
		if !ok {
			WriteError(c, errors.New("id must not be empty"))
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(c, err)
			return
		}
		buffer, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			WriteError(c, err)
			return
		}
		description = string(buffer)
		err = taskService.taskUsercase.UpdateDescription(c, int64(id), description)
		if err != nil {
			WriteError(c, err)
			return
		}
	})
	taskGroup.GET("/", func(c *gin.Context) {
		var minId int64
		tasks, err := taskService.taskUsercase.Get(c, minId, 10)
		if err != nil {
			WriteError(c, err)
			return
		}
		c.JSON(200, tasks)
	})
	taskGroup.PUT("/:id/complete", func(c *gin.Context) {
		idStr, ok := c.Params.Get("id")
		if !ok {
			WriteError(c, errors.New("id must not be empty"))
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			WriteError(c, err)
			return
		}
		err = taskService.taskUsercase.Complete(c, int64(id))
		if err != nil {
			WriteError(c, err)
			return
		}
	})

	return taskGroup
}
func WriteError(c *gin.Context, err error) {
	c.Data(500, "text/html", []byte(err.Error()))
	log.Println(err.Error())
}
