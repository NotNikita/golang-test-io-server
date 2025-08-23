package handlers

import (
	"context"

	"test-server/internal/domain/model"
)

//go:generate minimock -i TasksService -o ./mock -s _mock.go
type TasksService interface {
	RegisterTask(ctx context.Context, title string) (string, error)
	TaskInfo(ctx context.Context, taskId string) (*model.Task, error)
	DeleteTask(ctx context.Context, taskId string) error
}

type Handler struct {
	tasksService TasksService
}

func NewHandler(tasksService TasksService) *Handler {
	return &Handler{
		tasksService: tasksService,
	}
}
