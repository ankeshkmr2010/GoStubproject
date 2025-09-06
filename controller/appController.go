package controller

import (
	"github.com/gin-gonic/gin"
	"gostubproject/stub/interfaces"
)

type AppController struct {
	TaskDbAccessor interfaces.TasksWrapper
	// new repo proxy ork
}

func NewAppController(taskDbAccessor interfaces.TasksWrapper) *AppController {
	return &AppController{
		TaskDbAccessor: taskDbAccessor,
		// repo proxy
	}
}

func (a *AppController) Ping(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "pong"})
}
