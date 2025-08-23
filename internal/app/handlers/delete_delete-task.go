package handlers

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"test-server/internal/domain/model"
)

func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	taskId := c.Params("id")
	if !validateTaskId(taskId) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "error: task id is empty or has incorrect format",
		})
	}

	err := h.tasksService.DeleteTask(c.UserContext(), taskId)
	if err != nil {
		if errors.Is(err, model.ErrTaskNotFound) {
			return c.Status(fiber.StatusPreconditionFailed).JSON(fiber.Map{
			"ok":    false,
			"error": fmt.Errorf("task with provided id wasn't found: %w", err).Error(),
		})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": fmt.Errorf("failed to find task with provided id: %w", err).Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"ok":    true,
			"error": "task was successfully removed",
		})
}