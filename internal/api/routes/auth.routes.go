package routes

import (
	"food-ordering/internal/api/handlers"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(app fiber.Router, h *handlers.AuthHandler) {
	auth := app.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
}
