package handlers

import (
	"food-ordering/internal/dto"
	"food-ordering/internal/middleware"
	"food-ordering/internal/promo"
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderHandler struct {
	service             *services.OrderService
	idempotencyKeyIssue fiber.Handler // issuing handler built from middleware
}

func NewOrderHandler(service *services.OrderService, redisClient *redis.Client) *OrderHandler {
	return &OrderHandler{
		service:             service,
		idempotencyKeyIssue: middleware.IdempotencyKey(redisClient),
	}
}

func InitOrderModule(
	orderCollection *mongo.Collection,
	productCollection *mongo.Collection,
	promoValidator *promo.Validator,
	redisClient *redis.Client,
) *OrderHandler {
	orderRepo := repository.NewOrderRepository(orderCollection)
	productRepo := repository.NewProductRepository(productCollection)
	service := services.NewOrderService(orderRepo, productRepo, promoValidator)
	return NewOrderHandler(service, redisClient)
}

// GetIdempotencyKey godoc
// @Summary     Get an idempotency key
// @Description Generates a server-issued idempotency key scoped to the authenticated user. The key must be sent as the Idempotency-Key header when placing an order. Keys expire after 10 minutes if unused.
// @Tags        order
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  dto.IdempotencyKeyResponse
// @Failure     401  {object}  dto.ErrorResponse
// @Failure     500  {object}  dto.ErrorResponse
// @Header      200  {string}  Idempotency-Key  "The generated key"
// @Router      /order/idempotency-key [get]
func (h *OrderHandler) GetIdempotencyKey(ctx fiber.Ctx) error {
	return h.idempotencyKeyIssue(ctx)
}

// PlaceOrder godoc
// @Summary     Place an order
// @Description Places a new order for the authenticated user. Requires a server-issued idempotency key obtained from GET /order/idempotency-key. Sending the same key twice replays the original response without creating a duplicate order.
// @Tags        order
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       Idempotency-Key  header    string                 true  "Server-issued idempotency key (from GET /order/idempotency-key)"
// @Param       request          body      dto.PlaceOrderRequest  true  "Order details"
// @Success     200              {object}  dto.OrderSuccessResponse
// @Header      200              {string}  Idempotency-Key      "Echoed key"
// @Header      200              {string}  X-Idempotent-Replayed  "true if this is a replayed response"
// @Failure     400              {object}  dto.ErrorResponse
// @Failure     401              {object}  dto.ErrorResponse
// @Failure     409              {object}  dto.ErrorResponse  "Key already in use by a concurrent request"
// @Failure     422              {object}  dto.ErrorResponse  "Invalid/expired key or order validation error"
// @Router      /order [post]
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
