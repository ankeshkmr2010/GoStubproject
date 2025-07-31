package interfaces

import (
	"context"
	"github.com/google/uuid"
	"gostubproject/stub/models/dtos"
	s "gostubproject/stub/models/schemas"
)

type TaskDbAccessor interface {
	GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error)
	GetTasksForCreator(ctx context.Context, createdBy uuid.UUID) ([]dtos.GetTaskResp, error)
	CreateTask(ctx context.Context, task dtos.CreateTaskReq) (dtos.CreateTaskResp, error)
	UpdateTask(ctx context.Context, task dtos.UpdateTaskReq) (dtos.CreateTaskResp, error)
	DeleteTask(ctx context.Context, taskID uuid.UUID) error
	ListTasks(ctx context.Context) ([]dtos.GetTaskResp, error)
	ListTasksPaginated(ctx context.Context, statusFilter string, cursor *s.TaskCursor, pageSize int) (dtos.ListTasksResp, error)
}
