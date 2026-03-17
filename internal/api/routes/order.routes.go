package routes

import (
	"food-ordering/internal/api/handlers"
	"food-ordering/internal/middleware"
	"food-ordering/internal/services"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

func OrderRoutes(app fiber.Router, h *handlers.OrderHandler, authService *services.AuthService, redisClient *redis.Client) {
	jwtAuth := middleware.JWTAuth(authService)
	idempotency := middleware.Idempotency(redisClient)

	order := app.Group("/order", jwtAuth)
	order.Get("/idempotency-key", h.GetIdempotencyKey)
	order.Post("/", idempotency, h.PlaceOrder)
}
