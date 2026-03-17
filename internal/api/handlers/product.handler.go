package handlers

import (
	"food-ordering/internal/dto"
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

// GetProducts godoc
// @Summary     List all products
// @Description Returns all products available for ordering
// @Tags        product
// @Produce     json
// @Success     200  {object}  dto.ProductListResponse
// @Failure     500  {object}  dto.ErrorResponse
// @Router      /product [get]
func (h *ProductHandler) GetProducts(ctx fiber.Ctx) error {
	var products []dto.ProductResponse
	var err error

	products, err = h.service.GetProducts(ctx.Context())
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SuccessResponse(ctx, products)
}

// GetProduct godoc
// @Summary     Get product by ID
// @Description Returns a single product by its MongoDB ObjectID
// @Tags        product
// @Produce     json
// @Param       id   path      string  true  "Product ObjectID"
// @Success     200  {object}  dto.ProductSuccessResponse
// @Failure     400  {object}  dto.ErrorResponse  "Invalid ID"
// @Failure     404  {object}  dto.ErrorResponse  "Product not found"
// @Router      /product/{id} [get]
func (h *ProductHandler) GetProduct(ctx fiber.Ctx) error {
	id := ctx.Params("id")

	var product *dto.ProductResponse
	var err error

	product, err = h.service.GetProduct(ctx.Context(), id)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusNotFound, err.Error())
	}
	return utils.SuccessResponse(ctx, product)
}
