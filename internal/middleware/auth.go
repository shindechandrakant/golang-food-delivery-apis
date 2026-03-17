package middleware

import (
	"food-ordering/config"

	"github.com/gofiber/fiber/v3"
)

func ApiKeyAuth(ctx fiber.Ctx) error {
	apiKey := ctx.Get("api_key")
	expected := config.GetEnv("API_KEY")
	if expected == "" {
		expected = "apitest"
	}

	if apiKey != expected {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}
	return ctx.Next()
}
