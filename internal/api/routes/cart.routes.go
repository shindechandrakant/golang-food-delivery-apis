package routes

import (
	"food-ordering/internal/api/handlers"

	"github.com/gofiber/fiber/v3"
)

func CartRoutes(app fiber.Router, h *handlers.CartHandler) {

	cart := app.Group("/cart")
	cart.Post("/items", h.AddItem)
	cart.Get("/", h.GetCart)
	cart.Put("/items/:productId", h.UpdateItem)
	cart.Delete("/items/:productId", h.RemoveItem)
	cart.Delete("/", h.ClearCart)
}
