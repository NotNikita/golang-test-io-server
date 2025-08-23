package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"test-server/internal/domain/model"
	"time"

	"github.com/google/uuid"
)

type TasksRepository interface {
	CreateTask(ctx context.Context, task model.Task) error
	GetTask(ctx context.Context, id string) (*model.Task, error)
	UpdateTask(ctx context.Context, id string, status model.Status) error
	DeleteTask(ctx context.Context, id string) error
}

type TasksService struct {
	SaveInterval int // seconds
	tasksRepo    TasksRepository
}

func NewTasksService(interval int, tasksRepo TasksRepository) *TasksService {
	return &TasksService{
		SaveInterval: interval,
		tasksRepo:    tasksRepo,
	}
}

func (s *TasksService) RegisterTask(ctx context.Context, title string) (string, error) {
	task := model.Task{
		ID:        uuid.New(),
		Status:    model.Pending,
		Title:     title,
		CreatedAt: time.Now(),
	}

	err := s.tasksRepo.CreateTask(ctx, task)
	if err != nil {
		return "", fmt.Errorf("TasksService.RegisterTask: failed to create new task: %w", err)
	}

	go func() {
		sleepInterval := time.Duration(2+s.SaveInterval) * time.Second
		time.Sleep(sleepInterval)

		if rand.Float32() < 0.2 {
			task.Status = model.Failed
		} else {
			task.Status = model.Completed
		}

		task.Duration = sleepInterval
		if err := s.tasksRepo.UpdateTask(context.Background(), task.ID.String(), task.Status); err != nil {
			log.Printf("TasksService.RegisterTask: error while updating task status: %v", err.Error())
		}
	}()

	return task.ID.String(), nil
}

func (s *TasksService) TaskInfo(ctx context.Context, taskId string) (*model.Task, error) {
	taskInfo, err := s.tasksRepo.GetTask(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("TasksRepo.GetTask: failed to get task info by id: %w", err)
	}

	return taskInfo, nil
}

func (s *TasksService) DeleteTask(ctx context.Context, taskId string) error {
	err := s.tasksRepo.DeleteTask(ctx, taskId)
	if err != nil {
		return fmt.Errorf("TasksRepo.DeleteTask: failed to delete task info by id: %w", err)
	}

	return nil
}
