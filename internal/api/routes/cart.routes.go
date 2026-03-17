package routes

import (
	"food-ordering/internal/api/handlers"
	"food-ordering/internal/middleware"
	"food-ordering/internal/services"

	"github.com/gofiber/fiber/v3"
)

func CartRoutes(app fiber.Router, h *handlers.CartHandler, authService *services.AuthService) {
	jwtAuth := middleware.JWTAuth(authService)

	cart := app.Group("/cart", jwtAuth)
	cart.Post("/items", h.AddItem)
	cart.Get("/", h.GetCart)
	cart.Put("/items/:productId", h.UpdateItem)
	cart.Delete("/items/:productId", h.RemoveItem)
	cart.Delete("/", h.ClearCart)
}
