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

// Register godoc
// @Summary     Register a new user
// @Description Creates a new user account and returns a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request  body      dto.RegisterRequest       true  "Registration details"
// @Success     201      {object}  dto.AuthSuccessResponse
// @Failure     400      {object}  dto.ErrorResponse
// @Router      /auth/register [post]
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

// Login godoc
// @Summary     Login
// @Description Authenticates a user and returns a JWT token
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request  body      dto.LoginRequest          true  "Login credentials"
// @Success     200      {object}  dto.AuthSuccessResponse
// @Failure     400      {object}  dto.ErrorResponse
// @Failure     401      {object}  dto.ErrorResponse
// @Router      /auth/login [post]
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
