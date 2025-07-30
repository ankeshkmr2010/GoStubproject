package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gostubproject/stub/models/dtos"
)

func (a *AppController) GetTaskByID(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(400, gin.H{"error": "Task ID is required"})
		return
	}

	taskUuid, err := uuid.Parse(taskID)

	task, err := a.TaskDbAccessor.GetTaskByID(c.Request.Context(), taskUuid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve task"})
		return
	}

	c.JSON(200, task)
}

func (a *AppController) CreateTask(c *gin.Context) {
	var createTaskReq dtos.CreateTaskReq
	if err := c.ShouldBindJSON(&createTaskReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	createResp, err := a.TaskDbAccessor.CreateTask(c.Request.Context(), createTaskReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(201, createResp)
}

func (a *AppController) AddNewChildTask(c *gin.Context) {
	var addChildTaskReq dtos.AddNewChildTaskReq
	if err := c.ShouldBindJSON(&addChildTaskReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	err := a.TaskDbAccessor.AddNewChildTask(c.Request.Context(), addChildTaskReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add child task"})
		return
	}

	c.JSON(201, gin.H{"message": "Child task added successfully"})
}

func (a *AppController) UpdateTask(c *gin.Context) {
	var updateTaskReq dtos.UpdateTaskReq
	if err := c.ShouldBindJSON(&updateTaskReq); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	updateResp, err := a.TaskDbAccessor.UpdateTask(c.Request.Context(), updateTaskReq)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update task"})
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

	taskUuid, err := uuid.Parse(taskID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Task ID format"})
		return
	}

	err = a.TaskDbAccessor.DeleteTask(c.Request.Context(), taskUuid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(204, map[string]string{"message": "Task deleted successfully"})
}

func (a *AppController) ListAllTasks(c *gin.Context) {
	tasks, err := a.TaskDbAccessor.ListTasks(c.Request.Context())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	c.JSON(200, tasks)
}
