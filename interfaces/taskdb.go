package interfaces

import (
	"context"
	"github.com/google/uuid"
	"gostubproject/stub/models/dtos"
)

type TaskDbAccessor interface {
	GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error)
	GetTasksForCreator(ctx context.Context, createdBy uuid.UUID) ([]dtos.GetTaskResp, error)
	CreateTask(task dtos.CreateTaskReq) (dtos.CreateTaskResp, error)
	AddNewChildTask(task dtos.AddNewChildTaskReq) error

	UpdateTask(task dtos.UpdateTaskReq) (dtos.CreateTaskResp, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	DeleteTaskComment(ctx context.Context, taskID, commentID uuid.UUID) error
}
