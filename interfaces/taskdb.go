package interfaces

import (
	"context"
	"github.com/google/uuid"
	"gostubproject/stub/models/dtos"
)

type TaskDbAccessor interface {
	GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error)
	GetTasksForUserID(ctx context.Context, userID uuid.UUID) ([]dtos.GetTaskResp, error)
	CreateTask(task dtos.CreateTaskReq) (dtos.CreateTaskResp, error)

	UpdateTask(task dtos.UpdateTaskReq) (dtos.CreateTaskResp, error)
	UpdateTaskComment(comment dtos.UpdateTaskComment) (dtos.CreateTaskResp, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	DeleteTaskComment(ctx context.Context, taskID, commentID uuid.UUID) error
}
