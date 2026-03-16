package handlers

import (
	"food-ordering/internal/dto"
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{
		service: service,
	}
}

func InitCartModule(collection *mongo.Collection) *CartHandler {
	repo := repository.NewCartRepository(collection)
	service := services.NewCartService(repo)
	handler := NewCartHandler(service)
	return handler
}

func (h *CartHandler) AddItem(ctx fiber.Ctx) error {

	userId := ctx.Get("x-user-id")
	var req dto.AddCartItemRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}

	err := h.service.AddItem(ctx.Context(), userId, req)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}

	return utils.SuccessResponse(ctx, "Item added")
}

func (h *CartHandler) GetCart(ctx fiber.Ctx) error {
	userId := ctx.Get("x-user-id")
	cart, err := h.service.GetCart(ctx.Context(), userId)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}

	var totalCartValue float64 = 0
	for _, item := range cart.Items {
		totalCartValue += item.Price
	}
	return utils.SuccessResponse(ctx, fiber.Map{
		"cart":           cart,
		"totalCartValue": totalCartValue,
	})
}

func (h *CartHandler) ClearCart(ctx fiber.Ctx) error {
	userId := ctx.Get("x-user-id")

	err := h.service.ClearCart(ctx, userId)
	if err != nil {
		return utils.ErrorResponse(ctx, err.Error())
	}

	return utils.SuccessResponse(ctx, "Cart Cleared")
}
