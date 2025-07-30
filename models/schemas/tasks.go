package schemas

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Status      string          `json:"status"`
	Priority    int             `json:"priority"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	IsDeleted   bool            `json:"is_deleted"`
	TaskData    json.RawMessage `json:"task_data"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Children    []*Task         `json:"children,omitempty"`
}
