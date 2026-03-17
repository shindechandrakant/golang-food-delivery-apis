package handlers

import (
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(s *services.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func InitProductModule(collection *mongo.Collection) *ProductHandler {
	repo := repository.NewProductRepository(collection)
	service := services.NewProductService(repo)
	return NewProductHandler(service)
}

func (h *ProductHandler) GetProducts(ctx fiber.Ctx) error {

	products, err := h.service.GetProducts(ctx.Context())
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(ctx, products)
}

func (h *ProductHandler) GetProduct(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	product, err := h.service.GetProduct(ctx.Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusNotFound, err.Error())
	}
	return utils.SuccessResponse(ctx, product)
}
