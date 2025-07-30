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
	getLatestActiveTaskById = ` SELECT id, name, description, status, priority, created_by,is_deleted, task_data,created_at,updated_at
		FROM tasks
		WHERE id = $1
		  AND is_deleted = false
		LIMIT 1;`

	getLatestActiveTasksForCreator = ` SELECT id, name, description, status, priority, created_by, is_deleted,task_data,created_at,updated_at
	FROM tasks
	WHERE created_by = $1
	  AND is_deleted = false
	ORDER BY created_at DESC;
	`

	insertIntoTasks = `INSERT INTO tasks (id, name, description, status, priority, created_by, is_deleted, task_data, created_at, updated_at)
    	VALUES ($1, $2, $3, $4, $5, $6, false, $7, NOW(), NOW()) `

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

	listTasks = `SELECT id, name, description, status, priority, created_by, is_deleted, task_data, created_at, updated_at 
	FROM tasks
	WHERE is_deleted = false
	ORDER BY created_at DESC;`
)

type TasksDbAccessorImpl struct {
	db *pgxpool.Pool
}

func (t *TasksDbAccessorImpl) GetTaskByID(ctx context.Context, id uuid.UUID) (dtos.GetTaskResp, error) {
	var taskData s.Task
	row := t.db.QueryRow(ctx, getLatestActiveTaskById, id)
	err := row.Scan(&taskData.ID, &taskData.Name, &taskData.Description, &taskData.Status, &taskData.Priority, &taskData.CreatedBy, &taskData.IsDeleted, &taskData.TaskData, &taskData.CreatedAt, &taskData.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("No task found with ID: ", id)
			return dtos.GetTaskResp{}, nil
		}
		return dtos.GetTaskResp{}, err
	}

	return dtos.GetTaskResp{
		NonDeletedTaskResp: &dtos.NonDeletedTaskResp{
			ID:          taskData.ID,
			Name:        taskData.Name,
			Description: taskData.Description,
			Status:      taskData.Status,
			Priority:    taskData.Priority,
			CreatedBy:   taskData.CreatedBy,
			TaskData:    taskData.TaskData,
			CreatedAt:   taskData.CreatedAt,
			UpdatedAt:   taskData.UpdatedAt,
			Children:    nil, // Assuming children are not fetched here, can be added later
		},
	}, nil
}

func (t *TasksDbAccessorImpl) GetTasksForCreator(ctx context.Context, createdBy uuid.UUID) ([]dtos.GetTaskResp, error) {
	var tasks []dtos.GetTaskResp
	rows, err := t.db.Query(ctx, getLatestActiveTasksForCreator, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var taskData s.Task
		if err := rows.Scan(&taskData.ID, &taskData.Name, &taskData.Description, &taskData.Status, &taskData.Priority, &taskData.CreatedBy, &taskData.IsDeleted, &taskData.TaskData, &taskData.CreatedAt, &taskData.UpdatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, dtos.GetTaskResp{
			NonDeletedTaskResp: &dtos.NonDeletedTaskResp{
				ID:          taskData.ID,
				Name:        taskData.Name,
				Description: taskData.Description,
				Status:      taskData.Status,
				Priority:    taskData.Priority,
				CreatedBy:   taskData.CreatedBy,
				TaskData:    taskData.TaskData,
				CreatedAt:   taskData.CreatedAt,
				UpdatedAt:   taskData.UpdatedAt,
				Children:    nil, // Assuming children are not fetched here, can be added later
			},
			IsDeleted: taskData.IsDeleted,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (t *TasksDbAccessorImpl) CreateTask(ctx context.Context, task dtos.CreateTaskReq) (dtos.CreateTaskResp, error) {
	id := uuid.New()
	_, err := t.db.Exec(ctx, insertIntoTasks, id, task.Name, task.Description, task.Status, task.Priority, task.CreatedBy, task.TaskData)
	if err != nil {
		log.Println("Error inserting new task: ", err)
		return dtos.CreateTaskResp{}, err
	}

	return dtos.CreateTaskResp{
		ID: id,
		GetTaskResp: &dtos.GetTaskResp{
			NonDeletedTaskResp: &dtos.NonDeletedTaskResp{
				ID:          id,
				Name:        task.Name,
				Description: task.Description,
				Status:      task.Status,
				Priority:    task.Priority,
				CreatedBy:   task.CreatedBy,
				TaskData:    task.TaskData,
			},
		},
	}, nil
}

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

func (t *TasksDbAccessorImpl) UpdateTask(ctx context.Context, task dtos.UpdateTaskReq) (dtos.CreateTaskResp, error) {
	// get task by ID to ensure it exists and is not deleted
	existingTask, err := t.GetTaskByID(ctx, task.TaskID)
	if err != nil {
		log.Println("Error retrieving task for update: ", err)
		return dtos.CreateTaskResp{}, err
	}

	if existingTask.NonDeletedTaskResp == nil {
		log.Println("Task not found or is deleted: ", task.TaskID)
		return dtos.CreateTaskResp{}, fmt.Errorf("task not found or is deleted: %s", task.TaskID)
	}

	if existingTask.IsDeleted {
		log.Println("Task is deleted, cannot update: ", task.TaskID)
		return dtos.CreateTaskResp{}, fmt.Errorf("task is deleted: %s", task.TaskID)
	}

	if task.Name == "" {
		task.Name = existingTask.Name
	}
	if task.Description == "" {
		task.Description = existingTask.Description
	}
	if task.Status == "" {
		task.Status = existingTask.Status
	}
	if task.Priority == 0 {
		task.Priority = existingTask.Priority
	}
	if task.TaskData == nil {
		task.TaskData = existingTask.TaskData
	}

	tx, err := t.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return dtos.CreateTaskResp{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println("Rollback failed: ", err)
		}
	}(tx, ctx)

	_, err = tx.Exec(ctx, updateTask, task.Name, task.Description, task.Status, task.Priority, task.TaskData, task.TaskID)
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
		ID: task.TaskID,
		GetTaskResp: &dtos.GetTaskResp{
			NonDeletedTaskResp: &dtos.NonDeletedTaskResp{
				ID:          task.TaskID,
				Name:        task.Name,
				Description: task.Description,
				Status:      task.Status,
				Priority:    task.Priority,
				TaskData:    task.TaskData,
			},
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
		var taskData s.Task
		if err := rows.Scan(&taskData.ID, &taskData.Name, &taskData.Description, &taskData.Status, &taskData.Priority, &taskData.CreatedBy, &taskData.IsDeleted, &taskData.TaskData, &taskData.CreatedAt, &taskData.UpdatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, dtos.GetTaskResp{
			NonDeletedTaskResp: &dtos.NonDeletedTaskResp{
				ID:          taskData.ID,
				Name:        taskData.Name,
				Description: taskData.Description,
				Status:      taskData.Status,
				Priority:    taskData.Priority,
				CreatedBy:   taskData.CreatedBy,
				TaskData:    taskData.TaskData,
				CreatedAt:   taskData.CreatedAt,
				UpdatedAt:   taskData.UpdatedAt,
				Children:    nil, // Assuming children are not fetched here, can be added later
			},
			IsDeleted: taskData.IsDeleted,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func NewTaskDbAccessorImpl(db *pgxpool.Pool) interfaces.TaskDbAccessor {
	return &TasksDbAccessorImpl{
		db: db,
	}
}
