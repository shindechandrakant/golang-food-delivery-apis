package routes

import (
	"food-ordering/internal/api/handlers"
	"food-ordering/internal/middleware"
	"food-ordering/internal/services"

	"github.com/gofiber/fiber/v3"
)

func OrderRoutes(app fiber.Router, h *handlers.OrderHandler, authService *services.AuthService) {
	jwtAuth := middleware.JWTAuth(authService)

	order := app.Group("/order", jwtAuth)
	order.Post("/", h.PlaceOrder)
}
