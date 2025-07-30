package dtos

import (
	"encoding/json"
	"github.com/google/uuid"
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
	Status      string          `json:"status"`
	TaskData    json.RawMessage `json:"task_data"`
	Priority    int             `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
}

// --- Update Task ---

type UpdateTaskReq struct {
	TaskID      uuid.UUID       `json:"task_id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Status      string          `json:"status"`
	TaskData    json.RawMessage `json:"task_data"`
	Priority    int             `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
}
type AddNewChildTaskReq struct {
	TaskID    uuid.UUID   `json:"task_id"`
	ChildTask []uuid.UUID `json:"child_task"`
}

// --- Update Comment ---

// ---------------------------------------------------------------------------------
// --------------------------------- Response DTOs ---------------------------------
// ---------------------------------------------------------------------------------

// --- Get Task ---

type GetTaskResp struct {
	*NonDeletedTaskResp
	IsDeleted bool `json:"is_deleted"`
}

type NonDeletedTaskResp struct {
	ID          uuid.UUID       `json:"id"`
	Version     int             `json:"version"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Status      string          `json:"status"`
	Priority    int             `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	TaskData    json.RawMessage `json:"task_data"`
	CreatedAt   time.Time       `json:"created_at"`         // ISO 8601
	UpdatedAt   time.Time       `json:"updated_at"`         // ISO 8601
	Children    []*GetTaskResp  `json:"children,omitempty"` // Recursive structure for child
}

type CreateTaskResp struct {
	ID uuid.UUID `json:"id"`
	*GetTaskResp
}
