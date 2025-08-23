package repository

import (
	"context"
	"test-server/internal/domain/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGetTask(b *testing.B) {
	str := "2025-08-23T18:56:28.34065+02:00"
	timestamp, _ := time.Parse(time.RFC3339, str)
	testTaskInfo := model.Task{
		ID:        uuid.New(),
		Status:    model.Completed,
		Title:     "dummy-title",
		CreatedAt: timestamp,
		Duration:  time.Second * 3,
	}
	ctx := context.Background()

	repo := NewTasksRepository()
	err := repo.CreateTask(ctx, testTaskInfo)
	assert.NoError(b, err)
	
	taskID := testTaskInfo.ID.String()
	
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		obj, err := repo.GetTask(ctx, taskID)
		assert.NoError(b, err)
		assert.Equal(b, testTaskInfo.ID, obj.ID)
	}
}

func BenchmarkGetTask_Parallel(b *testing.B) {
		str := "2025-08-23T18:56:28.34065+02:00"
	timestamp, _ := time.Parse(time.RFC3339, str)
	testTaskInfo := model.Task{
		ID:        uuid.New(),
		Status:    model.Completed,
		Title:     "dummy-title",
		CreatedAt: timestamp,
		Duration:  time.Second * 3,
	}
	ctx := context.Background()

	repo := NewTasksRepository()
	err := repo.CreateTask(ctx, testTaskInfo)
	assert.NoError(b, err)
	
	taskID := testTaskInfo.ID.String()
	
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.GetTask(ctx, taskID)
			assert.NoError(b, err)
		}
	})
}