package middleware

import (
	"food-ordering/internal/services"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// JWTAuth returns a middleware that validates a Bearer token and sets
// "userId", "userEmail", and "userRole" in the request context locals.
func JWTAuth(authService *services.AuthService) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		header := ctx.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "missing or invalid authorization header",
			})
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := authService.ParseToken(tokenStr)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error":   "invalid or expired token",
			})
		}

		ctx.Locals("userId", claims.UserId)
		ctx.Locals("userEmail", claims.Email)
		ctx.Locals("userRole", claims.Role)

		return ctx.Next()
	}
}

// RequireRole returns a middleware that allows only users with one of the given roles.
// Must be used after JWTAuth.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(ctx fiber.Ctx) error {
		role, _ := ctx.Locals("userRole").(string)
		if _, ok := allowed[role]; !ok {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "forbidden",
			})
		}
		return ctx.Next()
	}
}
