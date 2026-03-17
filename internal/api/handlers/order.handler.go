package handlers

import (
	"food-ordering/internal/dto"
	"food-ordering/internal/promo"
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderHandler struct {
	service *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func InitOrderModule(
	orderCollection *mongo.Collection,
	productCollection *mongo.Collection,
	promoValidator *promo.Validator,
) *OrderHandler {
	orderRepo := repository.NewOrderRepository(orderCollection)
	productRepo := repository.NewProductRepository(productCollection)
	service := services.NewOrderService(orderRepo, productRepo, promoValidator)
	return NewOrderHandler(service)
}

func (h *OrderHandler) PlaceOrder(ctx fiber.Ctx) error {
	var req dto.PlaceOrderRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	if len(req.Items) == 0 {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, "items must not be empty")
	}

	order, err := h.service.PlaceOrder(ctx.Context(), req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnprocessableEntity, err.Error())
	}

	return utils.SuccessResponse(ctx, order)
}
