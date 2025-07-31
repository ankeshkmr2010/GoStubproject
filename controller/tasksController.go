package controller

import (
	"github.com/gin-gonic/gin"
	"gostubproject/stub/models/dtos"
	"log"
)

func (a *AppController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	task, err := a.TaskDbAccessor.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve task", "details": err.Error()})
		return
	}

	c.JSON(200, task)
}

func (a *AppController) CreateTask(c *gin.Context) {
	var createTaskReq dtos.CreateTaskReq
	if err := c.ShouldBindJSON(&createTaskReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	createResp, err := a.TaskDbAccessor.CreateTask(c.Request.Context(), createTaskReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task", "details": err.Error()})
		return
	}

	c.JSON(201, createResp)
}

func (a *AppController) UpdateTask(c *gin.Context) {
	var updateTaskReq dtos.UpdateTaskReq
	if err := c.ShouldBindJSON(&updateTaskReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	updateResp, err := a.TaskDbAccessor.UpdateTask(c.Request.Context(), updateTaskReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update task", "details": err.Error()})
		return
	}

	c.JSON(200, updateResp)
}

func (a *AppController) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	err := a.TaskDbAccessor.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete task", "details": err.Error()})
		return
	}

	c.JSON(204, gin.H{"message": "Task deleted successfully"})
}

func (a *AppController) ListAllTasks(c *gin.Context) {
	tasks, err := a.TaskDbAccessor.ListTasks(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve tasks", "details": err.Error()})
		return
	}
	c.JSON(200, tasks)
}

func (a *AppController) ListTasksPaginated(c *gin.Context) {
	var statusFilter string
	if status := c.Query("status"); status != "" {
		statusFilter = status
	} else {
		statusFilter = "all"
	}
	var listReq dtos.ListTasksReq
	if err := c.BindJSON(&listReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}
	resp, err := a.TaskDbAccessor.ListTasksPaginated(c.Request.Context(), statusFilter, listReq.Cursor, listReq.PageSize)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Failed to retrieve tasks", "details": err.Error()})
		return
	}
	c.JSON(200, resp)
}
