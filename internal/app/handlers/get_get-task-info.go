package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"test-server/internal/domain/model"
)

type taskInfoResponse struct {
	ID        uuid.UUID `json:"task_id"`
    Status    string    `json:"status"`
    Title     string    `json:"title"`
    CreatedAt time.Time `json:"created_at"`
    Duration  int64     `json:"duration_ms"` // Convert to milliseconds for API
}

type getTaskInfoResponse struct {
	TaskDetails taskInfoResponse `json:"data"`
	Error string `json:"error"`
	OK bool `json:"ok"`
}

func (h *Handler) GetTaskInfo(c *fiber.Ctx) error {
	taskId := c.Params("id")
	if !validateTaskId(taskId) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "error: task id is empty or has incorrect format",
		})
	}

	taskInfo, err := h.tasksService.TaskInfo(c.UserContext(), taskId)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
			"ok":    false,
			"error": fmt.Errorf("Task with provided id wasn't found: %w", err).Error(),
		})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": fmt.Errorf("Failed to find task with provided id: %w", err).Error(),
		})
	}
	if taskInfo == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok": false,
			"error": "service returned nil task without error",
		})
	}

	taskDTO := mapTaskToDTO(taskInfo)
	response := getTaskInfoResponse{
		OK: true,
		Error: "",
		TaskDetails: taskDTO,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func validateTaskId(s string) bool {
	if s == "" {
		return false
	}

	if uuid.Validate(s) != nil {
		return false
	}

	return true
}

func mapTaskToDTO(task *model.Task) taskInfoResponse {
    return taskInfoResponse{
        ID:        task.ID,
        Status:    string(task.Status), // Convert enum to string
        Title:     task.Title,
        CreatedAt: task.CreatedAt,
        Duration:  task.Duration.Milliseconds(), // Convert to milliseconds
    }
}