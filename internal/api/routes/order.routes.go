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

	// Step 1: client requests a server-issued idempotency key.
	order.Get("/idempotency-key", middleware.IdempotencyKey(redisClient))

	// Step 2: client places the order using that key.
	order.Post("/", idempotency, h.PlaceOrder)
}
