package controller

import (
	"github.com/gin-gonic/gin"
	"gostubproject/stub/interfaces"
)

type AppController struct {
	TaskDbAccessor interfaces.TaskDbAccessor
}

func NewAppController(taskDbAccessor interfaces.TaskDbAccessor) *AppController {
	return &AppController{
		TaskDbAccessor: taskDbAccessor,
	}
}

func (a *AppController) Ping(c *gin.Context) {
	c.JSON(200, map[string]string{"message": "pong"})
}
