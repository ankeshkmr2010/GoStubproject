package dtos

import (
	"encoding/json"
	"github.com/google/uuid"
	s "gostubproject/stub/models/schemas"
	"time"
)

// Request and Response DTOs for Task Management

// ---------------------------------------------------------------------------------
// --------------------------------- Request DTOs ---------------------------------
// ---------------------------------------------------------------------------------

// --- Get Task ---
// --- task id in path param --

// --- Create Task ---

type CreateTaskReq struct {
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Status      s.TaskStatus    `json:"status"`
	TaskData    json.RawMessage `json:"task_data"`
	Priority    int             `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	RequestedAt time.Time       `json:"requested_at"`
}

// --- Update Task ---

type UpdateTaskReq struct {
	TaskID      uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Status      s.TaskStatus    `json:"status"`
	TaskData    json.RawMessage `json:"task_data"`
	Priority    int             `json:"priority"`
	RequestedAt time.Time       `json:"requested_at"`
}
type AddNewChildTaskReq struct {
	TaskID    uuid.UUID   `json:"task_id"`
	ChildTask []uuid.UUID `json:"child_task"`
}

type ListTasksReq struct {
	Cursor   *s.TaskCursor `json:"cursor,omitempty"` // Optional cursor for pagination
	PageSize int           `json:"page_size"`        // Number of tasks to return per page
}

// ---------------------------------------------------------------------------------
// --------------------------------- Response DTOs ---------------------------------
// ---------------------------------------------------------------------------------

// --- Get Task ---

type GetTaskResp struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Status       s.TaskStatus    `json:"status"`
	Priority     int             `json:"priority"`
	CreatedBy    uuid.UUID       `json:"created_by"`
	TaskData     json.RawMessage `json:"task_data"`
	CreatedAt    time.Time       `json:"created_at,omitempty"`
	UpdatedAt    time.Time       `json:"updated_at,omitempty"`
	SerialNumber int64           `json:"serial_number,omitempty"`
	RequestedAt  time.Time       `json:"requested_at,omitempty"`
	IsDeleted    bool            `json:"is_deleted,omitempty"`
}

type NonDeletedTaskResp struct {
}

type CreateTaskResp struct {
	ID uuid.UUID `json:"id"`
	*GetTaskResp
}

type ListTasksResp struct {
	Tasks  []GetTaskResp `json:"tasks"`
	Cursor *s.TaskCursor `json:"cursor,omitempty"` // Optional cursor for pagination
}
