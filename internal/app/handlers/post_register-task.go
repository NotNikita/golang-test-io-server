package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type postRegisterTask struct {
	Title string `json:"title"`
}

func (h *Handler) PostRegisterTask(c *fiber.Ctx) error {
	var postRegisterTask postRegisterTask

	if err := c.BodyParser(&postRegisterTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "Cannot parse JSON",
		})
	}
	if postRegisterTask.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ok":    false,
			"error": "request's body doesnt match schema",
		})
	}

	newID, err := h.tasksService.RegisterTask(c.UserContext(), postRegisterTask.Title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"ok":    false,
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"ok":   true,
		"data": newID,
	})
}
