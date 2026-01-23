package schemas

import (
	"encoding/json"
	"errors"
	"fmt"
)

type TaskStatus int

const (
	StatusPending TaskStatus = iota
	StatusStarted
	StatusCompleted
	StatusFailed
)

var statusToString = map[TaskStatus]string{
	StatusPending:   "pending",
	StatusStarted:   "started",
	StatusCompleted: "completed",
	StatusFailed:    "failed",
}

var stringToStatus = map[string]TaskStatus{
	"pending":   StatusPending,
	"started":   StatusStarted,
	"completed": StatusCompleted,
	"failed":    StatusFailed,
}

func TaskStatusFromString(status string) (TaskStatus, error) {
	if status == "" {
		return StatusPending, errors.New("status is empty")
	}
	if status, ok := stringToStatus[status]; ok {
		return status, nil
	}
	return StatusPending, fmt.Errorf("invalid status: %s", status)
}

// String returns the string representation
func (s TaskStatus) String() string {
	return statusToString[s]
}

// MarshalJSON for JSON encoding
func (s TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON for JSON decoding
func (s *TaskStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	status, ok := stringToStatus[str]
	if !ok {
		return fmt.Errorf("invalid status: %s", str)
	}
	*s = status
	return nil
}
