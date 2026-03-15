package routes

import (
	"food-ordering/internal/api/handlers"

	"github.com/gofiber/fiber/v3"
)

func ProductRoutes(app fiber.Router, handler *handlers.ProductHandler) {
	product := app.Group("/products")

	product.Get("/", handler.GetProducts)
	product.Get("/:id", handler.GetProduct)
}
