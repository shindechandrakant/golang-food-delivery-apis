package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

const (
	IdempotencyKeyHeader  = "Idempotency-Key"
	idempotencyTTLDone    = 24 * time.Hour
	idempotencyTTLPending = 10 * time.Minute
	processingTTL         = 30 * time.Second

	statusPending    = "pending"
	statusProcessing = "processing"
)

type cachedResponse struct {
	StatusCode int             `json:"statusCode"`
	Body       json.RawMessage `json:"body"`
}

// IdempotencyKey returns a new server-issued idempotency key scoped to the requesting user.
// The key is stored in Redis as "pending" and expires in 10 minutes if unused.
//
// Must be called after JWTAuth so ctx.Locals("userId") is set.
func IdempotencyKey(redisClient *redis.Client) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userId, _ := ctx.Locals("userId").(string)

		key := generateUUID()
		redisKey := idempotencyRedisKey(userId, key)

		if err := redisClient.Set(context.Background(), redisKey, statusPending, idempotencyTTLPending).Err(); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "failed to generate idempotency key",
			})
		}

		ctx.Set(IdempotencyKeyHeader, key)
		return ctx.JSON(fiber.Map{
			"success":        true,
			"idempotencyKey": key,
			"expiresIn":      "10 minutes",
		})
	}
}

// Idempotency validates a server-issued idempotency key and prevents duplicate order processing.
//
// Key lifecycle:
//
//	"pending"    → key was issued, order not yet placed — proceed
//	"processing" → another request is actively using this key  — 409 Conflict
//	<json>       → order was already completed               — replay response
//	missing      → key is invalid or expired                 — 422
func Idempotency(redisClient *redis.Client) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		iKey := ctx.Get(IdempotencyKeyHeader)
		if iKey == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   fmt.Sprintf("%s header is required. Request a key from GET /api/order/idempotency-key", IdempotencyKeyHeader),
			})
		}

		userId, _ := ctx.Locals("userId").(string)
		redisKey := idempotencyRedisKey(userId, iKey)

		existing, err := redisClient.Get(context.Background(), redisKey).Result()
		if err == redis.Nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"success": false,
				"error":   "idempotency key is invalid or expired. Request a new key from GET /api/order/idempotency-key",
			})
		}
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "failed to validate idempotency key",
			})
		}

		switch existing {
		case statusProcessing:
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"success": false,
				"error":   "a request with this idempotency key is already being processed",
			})

		case statusPending:
			// First use of this key — transition to processing.
			set, err := redisClient.SetXX(
				context.Background(), redisKey, statusProcessing, processingTTL,
			).Result()
			if err != nil || !set {
				return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
					"success": false,
					"error":   "a request with this idempotency key is already being processed",
				})
			}

		default:
			// Existing JSON response — replay it.
			var cached cachedResponse
			if err := json.Unmarshal([]byte(existing), &cached); err == nil {
				ctx.Set("X-Idempotent-Replayed", "true")
				ctx.Set(IdempotencyKeyHeader, iKey)
				return ctx.Status(cached.StatusCode).Send(cached.Body)
			}
		}

		// Run the actual handler.
		handlerErr := ctx.Next()

		statusCode := ctx.Response().StatusCode()

		if statusCode < 200 || statusCode >= 300 {
			// On failure, reset to "pending" so the client can retry with the same key.
			redisClient.Set(context.Background(), redisKey, statusPending, idempotencyTTLPending)
			return handlerErr
		}

		// Cache the successful response for 24 hours.
		rawBody := make([]byte, len(ctx.Response().Body()))
		copy(rawBody, ctx.Response().Body())

		payload, err := json.Marshal(cachedResponse{
			StatusCode: statusCode,
			Body:       json.RawMessage(rawBody),
		})
		if err == nil {
			redisClient.Set(context.Background(), redisKey, payload, idempotencyTTLDone)
		}

		// Echo the key back in the response header.
		ctx.Set(IdempotencyKeyHeader, iKey)

		return handlerErr
	}
}

func idempotencyRedisKey(userId, key string) string {
	return fmt.Sprintf("idempotency:%s:%s", userId, key)
}
