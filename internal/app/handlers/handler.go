package handlers

import (
	"context"
)

type TasksService interface {
	RegisterTask(ctx context.Context, title string) (string, error)
}

type Handler struct {
	tasksService TasksService
}

func NewHandler(tasksService TasksService) *Handler {
	return &Handler{
		tasksService: tasksService,
	}
}
