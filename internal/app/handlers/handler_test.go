package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	mocks "test-server/internal/app/handlers/mock"
	"test-server/internal/domain/model"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gojuno/minimock/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTasksHandler_PostRegisterTask(t *testing.T) {
	t.Parallel()

	testTaskName := "Dummy Task"

	testTable := []struct {
		name         string
		body         map[string]string
		mockSetup    func(mc *minimock.Controller) TasksService
		expectedCode int
		expectedBody map[string]interface{}
		wantErr      require.ErrorAssertionFunc
	}{
		{
			name: "success",
			body: map[string]string{"title": testTaskName},
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc).RegisterTaskMock.Expect(minimock.AnyContext, testTaskName).Return("test-id", nil)
			},
			expectedCode: 200,
			expectedBody: map[string]interface{}{
				"ok":   true,
				"data": "test-id",
			},
			wantErr: require.NoError,
		},
		{
			name: "invalid JSON body",
			body: map[string]string{"field": "test"},
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc)
			},
			expectedCode: 400,
			expectedBody: map[string]interface{}{
				"ok":    false,
				"error": "request's body doesnt match schema",
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			service := tt.mockSetup(mc)

			handler := NewHandler(service)

			app := fiber.New()
			app.Post("/tasks", handler.PostRegisterTask)

			// Prepare request body
			var bodyReader io.Reader
			if tt.body != nil {
				bodyBytes, err := json.Marshal(tt.body)
				require.NoError(t, err)
				bodyReader = bytes.NewReader(bodyBytes)
			}

			// Create HTTP request
			req := httptest.NewRequest("POST", "/tasks", bodyReader)
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			tt.wantErr(t, err)

			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Parse JSON response
			var responseBody map[string]any
			err = json.Unmarshal(bodyBytes, &responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestTasksHandler_GetTaskInfo(t *testing.T) {
	t.Parallel()

	testTaskId := "ca545e27-4e9b-4c95-b38b-d72069e33975"
	id, _ := uuid.Parse(testTaskId)
	str := "2025-08-23T18:56:28.34065+02:00"
	timestamp, _ := time.Parse(time.RFC3339, str)
	testTaskInfo := &model.Task{
		ID:        id,
		Status:    model.Completed,
		Title:     "dummy-title",
		CreatedAt: timestamp,
		Duration:  time.Second * 3,
	}

	testTable := []struct {
		name         string
		path         string
		mockSetup    func(mc *minimock.Controller) TasksService
		expectedCode int
		expectedBody map[string]interface{}
		wantErr      require.ErrorAssertionFunc
	}{
		{
			name: "success",
			path: testTaskId,
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc).TaskInfoMock.Expect(minimock.AnyContext, testTaskId).Return(testTaskInfo, nil)
			},
			expectedCode: 200,
			expectedBody: map[string]any{
				"ok":    true,
				"error": "",
				"data": map[string]any{
					"title":       "dummy-title",
					"task_id":     testTaskId,
					"status":      "completed",
					"duration_ms": float64(3000),
					"created_at":  str,
				},
			},
			wantErr: require.NoError,
		},
		{
			name: "invalid request's path param",
			path: "incorrect-path",
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc)
			},
			expectedCode: 400,
			expectedBody: map[string]interface{}{
				"ok":    false,
				"error": "error: task id is empty or has incorrect format",
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			service := tt.mockSetup(mc)

			handler := NewHandler(service)

			app := fiber.New()
			app.Get("/tasks/:id", handler.GetTaskInfo)

			// Create HTTP request
			req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s", "/tasks", tt.path), &bytes.Reader{})
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			tt.wantErr(t, err)

			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Parse JSON response
			var responseBody map[string]any
			err = json.Unmarshal(bodyBytes, &responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestTasksHandler_DeleteTask(t *testing.T) {
	t.Parallel()

	testTaskId := "ca545e27-4e9b-4c95-b38b-d72069e33975"

	testTable := []struct {
		name         string
		path         string
		mockSetup    func(mc *minimock.Controller) TasksService
		expectedCode int
		expectedBody map[string]interface{}
		wantErr      require.ErrorAssertionFunc
	}{
		{
			name: "success",
			path: testTaskId,
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc).DeleteTaskMock.Expect(minimock.AnyContext, testTaskId).Return(nil)
			},
			expectedCode: 200,
			expectedBody: map[string]interface{}{
				"ok":   true,
				"data": "task was successfully removed",
			},
			wantErr: require.NoError,
		},
		{
			name: "invalid request's path param",
			path: "incorrect-path",
			mockSetup: func(mc *minimock.Controller) TasksService {
				return mocks.NewTasksServiceMock(mc)
			},
			expectedCode: 400,
			expectedBody: map[string]interface{}{
				"ok":    false,
				"error": "error: task id is empty or has incorrect format",
			},
			wantErr: require.NoError,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mc := minimock.NewController(t)
			service := tt.mockSetup(mc)

			handler := NewHandler(service)

			app := fiber.New()
			app.Delete("/tasks/:id", handler.DeleteTask)

			// Create HTTP request
			req := httptest.NewRequest("DELETE", fmt.Sprintf("%s/%s", "/tasks", tt.path), &bytes.Reader{})
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := app.Test(req)
			tt.wantErr(t, err)

			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			// Parse JSON response
			var responseBody map[string]any
			err = json.Unmarshal(bodyBytes, &responseBody)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}
