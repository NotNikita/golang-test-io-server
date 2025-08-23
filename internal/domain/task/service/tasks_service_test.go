package service

import (
	"context"
	"errors"
	"test-server/internal/domain/model"
	mocks "test-server/internal/domain/task/service/mock"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksService_RegisterTask(t *testing.T) {
	t.Parallel()

	testTaskTitle := "Test Task"

	testTable := []struct {
		name      string
		title     string
		interval  int
		mockSetup func(mc *minimock.Controller) TasksRepository
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name:     "success",
			title:    testTaskTitle,
			interval: 3,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).CreateTaskMock.Set(func(ctx context.Context, task model.Task) error {
					// Verify the task has expected values
					assert.Equal(t, testTaskTitle, task.Title)
					assert.Equal(t, model.Pending, task.Status)
					assert.NotEqual(t, uuid.Nil, task.ID)
					assert.False(t, task.CreatedAt.IsZero())
					return nil
				})
			},
			wantErr: require.NoError,
		},
		{
			name:     "repository error",
			title:    testTaskTitle,
			interval: 3,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).CreateTaskMock.Set(func(ctx context.Context, task model.Task) error {
					// Verify the task has expected values
					assert.Equal(t, testTaskTitle, task.Title)
					assert.Equal(t, model.Pending, task.Status)
					assert.NotEqual(t, uuid.Nil, task.ID)
					assert.False(t, task.CreatedAt.IsZero())
					return errors.New("repository error")
				})
			},
			wantErr: require.Error,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			repo := tt.mockSetup(mc)

			service := NewTasksService(tt.interval, repo)

			taskID, err := service.RegisterTask(context.Background(), tt.title)
			tt.wantErr(t, err)

			if err == nil {
				assert.NotEmpty(t, taskID)
				// Validate UUID format
				_, parseErr := uuid.Parse(taskID)
				assert.NoError(t, parseErr)
			}
		})
	}
}

func TestTasksService_TaskInfo(t *testing.T) {
	t.Parallel()

	testTaskID := "ca545e27-4e9b-4c95-b38b-d72069e33975"
	testTask := &model.Task{
		ID:        uuid.MustParse(testTaskID),
		Status:    model.Completed,
		Title:     "Test Task",
		CreatedAt: time.Now(),
		Duration:  time.Second * 5,
	}

	testTable := []struct {
		name      string
		taskID    string
		mockSetup func(mc *minimock.Controller) TasksRepository
		expected  *model.Task
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name:   "success",
			taskID: testTaskID,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).GetTaskMock.Expect(minimock.AnyContext, testTaskID).Return(testTask, nil)
			},
			expected: testTask,
			wantErr:  require.NoError,
		},
		{
			name:   "repository error",
			taskID: testTaskID,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).GetTaskMock.Expect(minimock.AnyContext, testTaskID).Return(nil, errors.New("task not found"))
			},
			expected: nil,
			wantErr:  require.Error,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			repo := tt.mockSetup(mc)

			service := NewTasksService(3, repo)

			task, err := service.TaskInfo(context.Background(), tt.taskID)
			tt.wantErr(t, err)

			if err == nil {
				assert.Equal(t, tt.expected, task)
			} else {
				assert.Nil(t, task)
			}
		})
	}
}

func TestTasksService_DeleteTask(t *testing.T) {
	t.Parallel()

	testTaskID := "ca545e27-4e9b-4c95-b38b-d72069e33975"

	testTable := []struct {
		name      string
		taskID    string
		mockSetup func(mc *minimock.Controller) TasksRepository
		wantErr   require.ErrorAssertionFunc
	}{
		{
			name:   "success",
			taskID: testTaskID,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).DeleteTaskMock.Expect(minimock.AnyContext, testTaskID).Return(nil)
			},
			wantErr: require.NoError,
		},
		{
			name:   "repository error",
			taskID: testTaskID,
			mockSetup: func(mc *minimock.Controller) TasksRepository {
				return mocks.NewTasksRepositoryMock(mc).DeleteTaskMock.Expect(minimock.AnyContext, testTaskID).Return(errors.New("delete failed"))
			},
			wantErr: require.Error,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			repo := tt.mockSetup(mc)

			service := NewTasksService(3, repo)

			err := service.DeleteTask(context.Background(), tt.taskID)
			tt.wantErr(t, err)
		})
	}
}
