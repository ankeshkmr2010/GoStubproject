package dtos

import (
	"github.com/google/uuid"
	"time"
)

// Request and Response DTOs for Task Management

// ---------------------------------------------------------------------------------
// --------------------------------- Request DTOs ---------------------------------
// ---------------------------------------------------------------------------------

// --- Get Task ---

type GetTask struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// --- Create Task ---

type CreateTaskReq struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	DueDate     time.Time `json:"due_date,omitempty"`
	AssignedTo  uuid.UUID `json:"assigned_to"`
	CreatedBy   uuid.UUID `json:"created_by"`
}

// --- Update Task ---

type UpdateTaskReq struct {
	TaskID      uuid.UUID `json:"task_id"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"` // ISO 8601 format
	AssignedTo  uuid.UUID `json:"assigned_to,omitempty"`
	UpdatedBy   uuid.UUID `json:"updated_by"`
}

// --- Update Comment ---

type UpdateTaskComment struct {
	TaskID    uuid.UUID `json:"task_id"`
	CommentId uuid.UUID `json:"comment_id"`
	Comment   string    `json:"comment"`
	CreatedBy uuid.UUID `json:"created_by"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"` // ISO 8601 format
	UpdatedAt time.Time `json:"updated_at"` // ISO 8601 format
}

// ---------------------------------------------------------------------------------
// --------------------------------- Response DTOs ---------------------------------
// ---------------------------------------------------------------------------------

// --- Get Task ---

type GetTaskResp struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Status      string        `json:"status"`
	Priority    int           `json:"priority"`
	DueDate     *time.Time    `json:"due_date,omitempty"` // ISO 8601
	AssignedTo  uuid.UUID     `json:"assigned_to"`
	CreatedAt   time.Time     `json:"created_at"`         // ISO 8601
	UpdatedAt   time.Time     `json:"updated_at"`         // ISO 8601
	Children    []GetTaskResp `json:"children,omitempty"` // Recursive structure for subtasks
}

type GetTaskCommentResp struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	Comment   string    `json:"comment"`
	CreatedBy uuid.UUID `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"` // ISO 8601

}

type CreateTaskResp struct {
	ID uuid.UUID `json:"id"`
	*GetTaskResp
}
