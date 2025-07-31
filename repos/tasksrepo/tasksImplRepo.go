package tasksrepo

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gostubproject/stub/constants"
	"gostubproject/stub/interfaces"
	"gostubproject/stub/models/dtos"
	s "gostubproject/stub/models/schemas"
	"log"
	"math"
	"time"
)

type TasksWrapperRepo struct {
	taskDb interfaces.TaskDbAccessor
}

func NewTasksRepo(taskDb interfaces.TaskDbAccessor) interfaces.TasksWrapper {
	return &TasksWrapperRepo{
		taskDb: taskDb,
	}
}

func (r *TasksWrapperRepo) GetTaskByID(ctx context.Context, id string) (dtos.GetTaskResp, error) {
	taskID, err := uuid.Parse(id)
	if err != nil {
		return dtos.GetTaskResp{}, err
	}
	if taskID == uuid.Nil {
		return dtos.GetTaskResp{}, errors.New("task ID is required for retrieval")
	}
	return r.taskDb.GetTaskByID(ctx, taskID)
}

func (r *TasksWrapperRepo) GetTasksForCreator(ctx context.Context, createdBy string) ([]dtos.GetTaskResp, error) {
	creatorID, err := uuid.Parse(createdBy)
	if err != nil {
		return nil, err
	}
	return r.taskDb.GetTasksForCreator(ctx, creatorID)
}

func (r *TasksWrapperRepo) CreateTask(ctx context.Context, task dtos.CreateTaskReq) (dtos.CreateTaskResp, error) {
	if task.RequestedAt.IsZero() {
		return dtos.CreateTaskResp{}, errors.New("requested at is required for task creation")
	}
	if task.Name == "" {
		return dtos.CreateTaskResp{}, errors.New("task name is required for creation")
	}
	return r.taskDb.CreateTask(ctx, task)
}

func (r *TasksWrapperRepo) UpdateTask(ctx context.Context, updTask dtos.UpdateTaskReq) (dtos.CreateTaskResp, error) {
	if updTask.TaskID == uuid.Nil {
		log.Println("Task ID is required for update")
		return dtos.CreateTaskResp{}, errors.New("task ID is required for update")
	}

	if updTask.RequestedAt.IsZero() {
		log.Println("Requested at is required for update")
		return dtos.CreateTaskResp{}, errors.New("requested at is required for update")
	}

	return r.taskDb.UpdateTask(ctx, updTask)
}

func (r *TasksWrapperRepo) DeleteTask(ctx context.Context, taskID string) error {
	taskUuid, err := uuid.Parse(taskID)
	if err != nil {
		return err
	}
	return r.taskDb.DeleteTask(ctx, taskUuid)
}

func (r *TasksWrapperRepo) ListTasks(ctx context.Context) ([]dtos.GetTaskResp, error) {
	return r.taskDb.ListTasks(ctx)
}

func (r *TasksWrapperRepo) ListTasksPaginated(ctx context.Context, statusFilter string, cursor *s.TaskCursor, pageSize int) (dtos.ListTasksResp, error) {
	if cursor == nil || (cursor.CreatedAt.IsZero() && cursor.SerialNumber == 0) {
		log.Println("Cursor is empty")
		cursor.CreatedAt = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		cursor.SerialNumber = math.MinInt64
	}

	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}

	return r.taskDb.ListTasksPaginated(ctx, statusFilter, cursor, pageSize)
}
