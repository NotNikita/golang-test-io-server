package repository

import (
	"context"
	"sync"
	"test-server/internal/domain/model"
)

type TasksRepository struct {
	storage map[string]model.Task
	mu      sync.RWMutex
}

func NewTasksRepository() *TasksRepository {
	return &TasksRepository{
		storage: make(map[string]model.Task),
		mu:      sync.RWMutex{},
	}
}

func (repo *TasksRepository) CreateTask(ctx context.Context, task model.Task) error {
	id := task.ID.String()

	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.storage[id]
	if exists {
		return model.ErrTaskAlreadyExists
	}

	repo.storage[id] = task
	return nil
}

func (repo *TasksRepository) GetTask(ctx context.Context, id string) (*model.Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	task, exists := repo.storage[id]
	if !exists {
		return nil, model.ErrTaskNotFound
	}

	return &task, nil
}

func (repo *TasksRepository) UpdateTask(ctx context.Context, id string, status model.Status) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	
	task, exists := repo.storage[id]
	if !exists {
		return model.ErrTaskNotFound
	}
	
	task.Status = status
	repo.storage[id] = task
	return nil
}

func (repo *TasksRepository) DeleteTask(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	_, exists := repo.storage[id]
	if !exists {
		return model.ErrTaskNotFound
	}

	delete(repo.storage, id)
	return nil
}