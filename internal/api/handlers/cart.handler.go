package handlers

import (
	"food-ordering/internal/dto"
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CartHandler struct {
	service *services.CartService
}

func NewCartHandler(service *services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func InitCartModule(collection *mongo.Collection, redisClient *redis.Client) *CartHandler {
	mongoRepo := repository.NewMongoCartRepository(collection)
	redisRepo := repository.NewRedisCartRepository(redisClient)
	repo := repository.NewWriteThroughCartRepository(mongoRepo, redisRepo)
	service := services.NewCartService(repo)
	return NewCartHandler(service)
}

// AddItem godoc
// @Summary     Add item to cart
// @Description Adds a product to the authenticated user's cart. If the product is already in the cart, its quantity is incremented.
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       request  body      dto.AddCartItemRequest  true  "Product to add"
// @Success     200      {object}  dto.MessageResponse
// @Failure     400      {object}  dto.ErrorResponse
// @Failure     401      {object}  dto.ErrorResponse
// @Router      /cart/items [post]
func (h *CartHandler) AddItem(ctx fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(string)
	var req dto.AddCartItemRequest

	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	err := h.service.AddItem(ctx.Context(), userId, req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, "Item added")
}

// GetCart godoc
// @Summary     Get cart
// @Description Returns the authenticated user's cart with total value
// @Tags        cart
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  dto.CartSuccessResponse
// @Failure     401  {object}  dto.ErrorResponse
// @Failure     404  {object}  dto.ErrorResponse
// @Router      /cart [get]
func (h *CartHandler) GetCart(ctx fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(string)
	cart, err := h.service.GetCart(ctx.Context(), userId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusNotFound, err.Error())
	}

	var totalCartValue float64 = 0
	for _, item := range cart.Items {
		totalCartValue += item.Price * float64(item.Quantity)
	}
	return utils.SuccessResponse(ctx, fiber.Map{
		"cart":           cart,
		"totalCartValue": totalCartValue,
	})
}

// UpdateItem godoc
// @Summary     Update cart item quantity
// @Description Updates the quantity of a specific product in the cart
// @Tags        cart
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       productId  path      string                    true  "Product ObjectID"
// @Param       request    body      dto.UpdateCartItemRequest true  "New quantity"
// @Success     200        {object}  dto.MessageResponse
// @Failure     400        {object}  dto.ErrorResponse
// @Failure     401        {object}  dto.ErrorResponse
// @Router      /cart/items/{productId} [put]
func (h *CartHandler) UpdateItem(ctx fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(string)
	productId := ctx.Params("productId")

	var req dto.UpdateCartItemRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	err := h.service.UpdateItem(ctx.Context(), userId, productId, req.Quantity)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, "Item updated")
}

// RemoveItem godoc
// @Summary     Remove item from cart
// @Description Removes a specific product from the cart
// @Tags        cart
// @Produce     json
// @Security    BearerAuth
// @Param       productId  path      string  true  "Product ObjectID"
// @Success     200        {object}  dto.MessageResponse
// @Failure     400        {object}  dto.ErrorResponse
// @Failure     401        {object}  dto.ErrorResponse
// @Router      /cart/items/{productId} [delete]
func (h *CartHandler) RemoveItem(ctx fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(string)
	productId := ctx.Params("productId")

	err := h.service.RemoveItem(ctx.Context(), userId, productId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, "Item removed")
}

// ClearCart godoc
// @Summary     Clear cart
// @Description Removes all items from the authenticated user's cart
// @Tags        cart
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  dto.MessageResponse
// @Failure     400  {object}  dto.ErrorResponse
// @Failure     401  {object}  dto.ErrorResponse
// @Router      /cart [delete]
func (h *CartHandler) ClearCart(ctx fiber.Ctx) error {
	userId, _ := ctx.Locals("userId").(string)

	err := h.service.ClearCart(ctx, userId)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, "Cart Cleared")
}
