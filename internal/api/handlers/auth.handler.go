package handlers

import (
	"food-ordering/internal/dto"
	"food-ordering/internal/repository"
	"food-ordering/internal/services"
	"food-ordering/internal/utils"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func InitAuthModule(userCollection *mongo.Collection, jwtSecret string) (*AuthHandler, *services.AuthService) {
	repo := repository.NewUserRepository(userCollection)
	service := services.NewAuthService(repo, jwtSecret)
	return NewAuthHandler(service), service
}

func (h *AuthHandler) Register(ctx fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.service.Register(ctx.Context(), req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data":    resp,
	})
}

func (h *AuthHandler) Login(ctx fiber.Ctx) error {
	var req dto.LoginRequest
	if err := ctx.Bind().Body(&req); err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.service.Login(ctx.Context(), req)
	if err != nil {
		return utils.ErrorResponse(ctx, fiber.StatusUnauthorized, err.Error())
	}

	return utils.SuccessResponse(ctx, resp)
}
