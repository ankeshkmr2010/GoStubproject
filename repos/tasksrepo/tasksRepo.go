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
	"math"
	"time"
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

	addChildTask = `INSERT INTO task_relationships(parent_id, child_id, created_at) VALUES 
	($1, $2, NOW())`

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
	if id == uuid.Nil {
		log.Println("Task ID is required for retrieval")
		return dtos.GetTaskResp{}, errors.New("task ID is required for retrieval")
	}
	row := t.db.QueryRow(ctx, getLatestActiveTaskById, id)
	err := row.Scan(&taskData.ID, &taskData.Name, &taskData.Description, &taskData.Status, &taskData.Priority, &taskData.CreatedBy, &taskData.IsDeleted, &taskData.TaskData, &taskData.CreatedAt, &taskData.UpdatedAt, &taskData.RequestedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No task found with ID: ", id)
			return dtos.GetTaskResp{}, nil
		}
		return dtos.GetTaskResp{}, err
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
		if err := rows.Scan(&td.ID, &td.Name, &td.Description, &td.Status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, td.RequestedAt); err != nil {
			return nil, err
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
	if task.RequestedAt.IsZero() {
		log.Println("Requested at is required for task creation")
		return dtos.CreateTaskResp{}, errors.New("requested at is required for task creation")
	}
	if task.Name == "" {
		log.Println("Task name is required for creation")
		return dtos.CreateTaskResp{}, errors.New("task name is required for creation")
	}

	if task.Status == "" {
		log.Println("Task status is required for creation")
		return dtos.CreateTaskResp{}, errors.New("task status is required for creation")
	}

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

// AddNewChildTask : For later functionality, currently not used
func (t *TasksDbAccessorImpl) AddNewChildTask(ctx context.Context, task dtos.AddNewChildTaskReq) error {
	if len(task.ChildTask) == 0 {
		return errors.New("no child tasks provided")
	}

	//  change this to bulk inserts
	for _, childID := range task.ChildTask {
		_, err := t.db.Exec(ctx, addChildTask, task.TaskID, childID)
		if err != nil {
			log.Println("Error adding child task: ", err)
			return err
		}
	}

	return nil
}

func (t *TasksDbAccessorImpl) UpdateTask(ctx context.Context, updTask dtos.UpdateTaskReq) (dtos.CreateTaskResp, error) {
	// get task by ID to ensure it exists and is not deleted
	if updTask.TaskID == uuid.Nil {
		log.Println("Task ID is required for update")
		return dtos.CreateTaskResp{}, errors.New("task ID is required for update")
	}

	if updTask.RequestedAt.IsZero() {
		log.Println("Requested at is required for update")
		return dtos.CreateTaskResp{}, errors.New("requested at is required for update")
	}

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

	err = tx.QueryRow(ctx, getTaskDetailsForUpdate, updTask.TaskID).Scan(&et.ID, &et.SerialNumber, &et.Name, &et.Description, &et.Status, &et.Priority, &et.CreatedBy, &et.IsDeleted, &et.TaskData, &et.CreatedAt, &et.UpdatedAt, &et.RequestedAt)
	if err != nil {
		log.Println("Error retrieving task for update: ", err)
		return dtos.CreateTaskResp{}, err
	}

	if updTask.RequestedAt.Before(et.RequestedAt) {
		log.Println("RequestedAt for update is older than existing task")
		return dtos.CreateTaskResp{}, errors.New("requestedAt for update task is lesser than existing tasks requested at time")
	}

	if updTask.Name == "" {
		updTask.Name = et.Name
	}
	if updTask.Description == "" {
		updTask.Description = et.Description
	}
	if updTask.Status == "" {
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
		log.Println("Error deleting task: ", err)
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
		if err := rows.Scan(&td.ID, &td.Name, &td.Description, &td.Status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, &td.RequestedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, getRespFromTaskData(td))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TasksDbAccessorImpl) ListTasksPaginated(ctx context.Context, statusFilter string, cursor *s.TaskCursor, pageSize int) (dtos.ListTasksResp, error) {

	// Handle default cursor for first page
	var createdAt time.Time
	var serialNumber int64
	if cursor == nil || (cursor.CreatedAt.IsZero() && cursor.SerialNumber == 0) {
		log.Println("Cursor is empty")
		createdAt = time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		serialNumber = math.MinInt64
	} else {
		createdAt = cursor.CreatedAt
		serialNumber = cursor.SerialNumber
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	if statusFilter == "" || statusFilter == "all" {
		statusFilter = "%%"
	}
	query := listAllPaginated
	rows, err := t.db.Query(ctx, query, createdAt, serialNumber, pageSize, statusFilter)
	if err != nil {
		return dtos.ListTasksResp{}, err
	}
	defer rows.Close()

	tasks := make([]dtos.GetTaskResp, 0)
	var nextCursor *s.TaskCursor

	for rows.Next() {
		var td s.Task
		err := rows.Scan(&td.ID, &td.Name, &td.SerialNumber, &td.Description, &td.Status, &td.Priority, &td.CreatedBy, &td.IsDeleted, &td.TaskData, &td.CreatedAt, &td.UpdatedAt, &td.RequestedAt)
		if err != nil {
			return dtos.ListTasksResp{}, err
		}
		tasks = append(tasks, getRespFromTaskData(td))

		// Set next cursor to last item (assuming DESC order)
		nextCursor = &s.TaskCursor{
			CreatedAt:    td.CreatedAt,
			SerialNumber: td.SerialNumber,
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
