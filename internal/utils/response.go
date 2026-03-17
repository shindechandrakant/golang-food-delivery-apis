package utils

import "github.com/gofiber/fiber/v3"

func SuccessResponse(ctx fiber.Ctx, data interface{}) error {
	return ctx.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func ErrorResponse(ctx fiber.Ctx, status int, message string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
