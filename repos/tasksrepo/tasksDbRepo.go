package tasksrepo

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gostubproject/stub/interfaces"
	"gostubproject/stub/models/dtos"
	s "gostubproject/stub/models/schemas"
	"log"
)

const (
	getLatestActiveTaskById = ` SELECT id, name, description, status, priority, created_by,is_deleted, task_data,created_at,updated_at, requested_at
		FROM tasks
		WHERE id = $1
		  AND is_deleted = false
		LIMIT 1;`

	getLatestActiveTasksForCreator = ` SELECT id, name, description, status, priority, created_by, is_deleted,task_data,created_at,updated_at, requested_at
	FROM tasks
	WHERE created_by = $1
	  AND is_deleted = false
	ORDER BY created_at DESC;
	`

	getTaskDetailsForUpdate = `
	SELECT id, serial_number, name, description, status, priority,created_by, is_deleted, task_data, created_at, updated_at, requested_at
		FROM tasks
		WHERE id = $1 
		  AND is_deleted = false
		FOR UPDATE;`

	insertIntoTasks = `INSERT INTO tasks (id, name, description, status, priority, created_by, is_deleted, task_data, requested_at,created_at, updated_at)
    	VALUES ($1, $2, $3, $4, $5, $6, false, $7,$8, NOW(), NOW()) `

	updateTask = `UPDATE tasks
	SET name = $1, description = $2, status = $3, priority = $4, task_data = $5, updated_at = NOW()
	WHERE id = $6 and 
	    is_deleted = false`

	deleteTask = `UPDATE tasks
	SET is_deleted = true, updated_at = NOW()
	WHERE id = $1
	  AND is_deleted = false;`

	listTasks = `SELECT id, name, description, status, priority, created_by, is_deleted, task_data, created_at, updated_at, requested_at 
	FROM tasks
	WHERE is_deleted = false
	ORDER BY created_at DESC;`

	listAllPaginated = `
		SELECT id, name, serial_number, description, status, priority, created_by, is_deleted, task_data, created_at, updated_at, requested_at
		FROM tasks
		WHERE is_deleted = false
		  AND status ILIKE $4
		AND (
		    created_at > $1
		    OR (created_at = $1 AND serial_number > $2)
		)
		ORDER BY created_at , serial_number 
		LIMIT $3;
	`
)

type TasksDbAccessorImpl struct {
	db *pgxpool.Pool
}

func (t *TasksDbAccessorImpl) GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error) {
	var taskData s.Task
	row := t.db.QueryRow(ctx, getLatestActiveTaskById, id)
	var status string
	err := row.Scan(&taskData.ID, &taskData.Name, &taskData.Description, &status, &taskData.Priority, &taskData.CreatedBy, &taskData.IsDeleted, &taskData.TaskData, &taskData.CreatedAt, &taskData.UpdatedAt, &taskData.RequestedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No task found with ID: ", id)
			return dtos.GetTaskResp{}, nil
		}
		return dtos.GetTaskResp{}, err
	}
	taskData.Status, err = s.TaskStatusFromString(status)
	if err != nil {
		return dtos.GetTaskResp{}, fmt.Errorf("error converting status string to TaskStatus: %w", err)
	}

	return getRespFromTaskData(taskData), nil
}

func (t *TasksDbAccessorImpl) GetTasksForCreator(ctx context.Context, createdBy uuid.UUID) ([]dtos.GetTaskResp, error) {
	var tasks []dtos.GetTaskResp
	rows, err := t.db.Query(ctx, getLatestActiveTasksForCreator, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var td s.Task
		var status string
		if err := rows.Scan(&td.ID, &td.Name, &td.Description, &status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, td.RequestedAt); err != nil {
			return nil, err
		}
		td.Status, err = s.TaskStatusFromString(status)
		if err != nil {
			log.Println("Error converting status string to TaskStatus: ", err)
			return nil, fmt.Errorf("error converting status string to TaskStatus: %w", err)
		}
		tasks = append(tasks, getRespFromTaskData(td))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (t *TasksDbAccessorImpl) CreateTask(ctx context.Context, task dtos.CreateTaskReq) (dtos.CreateTaskResp, error) {
	id := uuid.New()

	_, err := t.db.Exec(ctx, insertIntoTasks, id, task.Name, task.Description, task.Status, task.Priority, task.CreatedBy, task.TaskData, task.RequestedAt)
	if err != nil {
		log.Println("Error inserting new task: ", err)
		return dtos.CreateTaskResp{}, err
	}

	return dtos.CreateTaskResp{
		ID: id,
		GetTaskResp: &dtos.GetTaskResp{
			ID:          id,
			Name:        task.Name,
			Description: task.Description,
			Status:      task.Status,
			Priority:    task.Priority,
			TaskData:    task.TaskData,
			CreatedBy:   task.CreatedBy,
			RequestedAt: task.RequestedAt,
		},
	}, nil
}

func (t *TasksDbAccessorImpl) UpdateTask(ctx context.Context, updTask dtos.UpdateTaskReq) (dtos.CreateTaskResp, error) {
	// get task by ID to ensure it exists and is not deleted

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return dtos.CreateTaskResp{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Println("Rollback failed: ", err)
		}
	}(tx, ctx)
	var et s.Task
	var status string
	err = tx.QueryRow(ctx, getTaskDetailsForUpdate, updTask.TaskID).Scan(&et.ID, &et.SerialNumber, &et.Name, &et.Description, &status, &et.Priority, &et.CreatedBy, &et.IsDeleted, &et.TaskData, &et.CreatedAt, &et.UpdatedAt, &et.RequestedAt)
	if err != nil {
		log.Println("Error retrieving task for update: ", err)
		return dtos.CreateTaskResp{}, err
	}
	et.Status, err = s.TaskStatusFromString(status)
	if err != nil {
		return dtos.CreateTaskResp{}, fmt.Errorf("error converting status string to TaskStatus: %w", err)
	}

	if updTask.RequestedAt.Before(et.RequestedAt) {
		return dtos.CreateTaskResp{}, errors.New("requestedAt for update task is lesser than existing tasks requested at time")
	}

	if updTask.Name == "" {
		updTask.Name = et.Name
	}
	if updTask.Description == "" {
		updTask.Description = et.Description
	}
	if updTask.Status == 0 {
		updTask.Status = et.Status
	}
	if updTask.Priority == 0 {
		updTask.Priority = et.Priority
	}
	if updTask.TaskData == nil {
		updTask.TaskData = et.TaskData
	}

	_, err = tx.Exec(ctx, updateTask, updTask.Name, updTask.Description, updTask.Status, updTask.Priority, updTask.TaskData, updTask.TaskID)
	if err != nil {
		log.Println("Error updating task: ", err)
		return dtos.CreateTaskResp{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error committing transaction: ", err)
		return dtos.CreateTaskResp{}, err
	}

	return dtos.CreateTaskResp{
		ID: updTask.TaskID,
		GetTaskResp: &dtos.GetTaskResp{
			ID:          updTask.TaskID,
			Name:        updTask.Name,
			Description: updTask.Description,
			Status:      updTask.Status,
			Priority:    updTask.Priority,
			TaskData:    updTask.TaskData,
		},
	}, nil
}

func (t *TasksDbAccessorImpl) DeleteTask(ctx context.Context, taskID uuid.UUID) error {
	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("Rollback failed: ", err)
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx, deleteTask, taskID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (t *TasksDbAccessorImpl) ListTasks(ctx context.Context) ([]dtos.GetTaskResp, error) {
	var tasks []dtos.GetTaskResp
	rows, err := t.db.Query(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var td s.Task
		var status string
		if err := rows.Scan(&td.ID, &td.Name, &td.Description, &status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, &td.RequestedAt); err != nil {
			return nil, err
		}
		td.Status, err = s.TaskStatusFromString(status)
		if err != nil {
			log.Println("Error converting status string to TaskStatus: ", err)
			return nil, fmt.Errorf("error converting status string to TaskStatus: %w", err)
		}

		tasks = append(tasks, getRespFromTaskData(td))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TasksDbAccessorImpl) ListTasksPaginated(ctx context.Context, statusFilter string, cursor *s.TaskCursor, pageSize int) (dtos.ListTasksResp, error) {
	if statusFilter == "" || statusFilter == "all" {
		statusFilter = "%%"
	} else {
		_, err := s.TaskStatusFromString(statusFilter)
		if err != nil {
			log.Println("Invalid status filter provided: ", statusFilter)
			return dtos.ListTasksResp{}, fmt.Errorf("invalid status filter: %s", statusFilter)
		}
	}
	rows, err := t.db.Query(ctx, listAllPaginated, cursor.CreatedAt, cursor.SerialNumber, pageSize, statusFilter)
	if err != nil {
		return dtos.ListTasksResp{}, err
	}
	defer rows.Close()

	tasks := make([]dtos.GetTaskResp, 0)
	var nextCursor *s.TaskCursor
	var lastTask *s.Task
	for rows.Next() {
		var td s.Task
		var status string
		err := rows.Scan(&td.ID, &td.Name, &td.SerialNumber, &td.Description, &status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, &td.RequestedAt)
		if err != nil {
			return dtos.ListTasksResp{}, err
		}
		td.Status, err = s.TaskStatusFromString(status)
		if err != nil {
			return dtos.ListTasksResp{}, fmt.Errorf("error converting status string to TaskStatus: %w", err)
		}
		tasks = append(tasks, getRespFromTaskData(td))
		lastTask = &td
	}

	if lastTask != nil {
		nextCursor = &s.TaskCursor{
			CreatedAt:    lastTask.CreatedAt,
			SerialNumber: lastTask.SerialNumber,
		}
	}
	return dtos.ListTasksResp{Tasks: tasks, Cursor: nextCursor}, rows.Err()
}

func getRespFromTaskData(taskData s.Task) dtos.GetTaskResp {
	return dtos.GetTaskResp{
		ID:          taskData.ID,
		Name:        taskData.Name,
		Description: taskData.Description,
		Status:      taskData.Status,
		Priority:    taskData.Priority,
		CreatedBy:   taskData.CreatedBy,
		TaskData:    taskData.TaskData,
		CreatedAt:   taskData.CreatedAt,
		UpdatedAt:   taskData.UpdatedAt,
		RequestedAt: taskData.RequestedAt,
		IsDeleted:   taskData.IsDeleted,
	}
}

func NewTaskDbAccessorImpl(db *pgxpool.Pool) interfaces.TaskDbAccessor {
	return &TasksDbAccessorImpl{
		db: db,
	}
}
