package model

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	Pending   Status = "pending"
	Completed Status = "completed"
	Failed    Status = "failed"
)

type Task struct {
	ID        uuid.UUID     `json:"task_id"`
	Status    Status        `json:"status"`
	Title     string        `json:"title"`
	CreatedAt time.Time     `json:"created_at"`
	Duration  time.Duration `json:"duration"`
}
