package interfaces

import (
	"context"
	"github.com/google/uuid"
	"gostubproject/stub/models/dtos"
)

type TaskDbAccessor interface {
	GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error)
	GetTasksForCreator(ctx context.Context, createdBy uuid.UUID) ([]dtos.GetTaskResp, error)
	CreateTask(ctx context.Context, task dtos.CreateTaskReq) (dtos.CreateTaskResp, error)
	AddNewChildTask(ctx context.Context, task dtos.AddNewChildTaskReq) error

	UpdateTask(ctx context.Context, task dtos.UpdateTaskReq) (dtos.CreateTaskResp, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	ListTasks(ctx context.Context) ([]dtos.GetTaskResp, error)
}
