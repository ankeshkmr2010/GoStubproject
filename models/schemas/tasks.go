package schemas

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	AssignedTo  uuid.UUID  `json:"assigned_to"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	IsDeleted   bool       `json:"is_deleted"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Children    []Task     `json:"children,omitempty"`
}

type TaskComment struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	Comment   string    `json:"comment"`
	CreatedBy uuid.UUID `json:"created_by"`
	IsDeleted bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
