package repository

import (
	"context"
	"errors"
	"food-ordering/internal/models"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// WriteThroughCartRepository implements the write-through caching pattern:
//
//   - Read:  Redis first → on cache miss, read from MongoDB → backfill Redis → return
//   - Write: MongoDB first (source of truth) → then Redis
//   - Clear: MongoDB first → then Redis
//
// MongoDB failures on write abort the operation.
// Redis failures on write are logged but do not abort — the DB is the source of truth.
type WriteThroughCartRepository struct {
	mongo *MongoCartRepository
	redis *RedisCartRepository
}

func NewWriteThroughCartRepository(
	mongo *MongoCartRepository,
	redis *RedisCartRepository,
) *WriteThroughCartRepository {
	return &WriteThroughCartRepository{mongo: mongo, redis: redis}
}

func (w *WriteThroughCartRepository) FindByUser(ctx context.Context, userId string) (*models.Cart, error) {
	// 1. Try Redis cache first.
	cart, err := w.redis.FindByUser(ctx, userId)
	if err == nil && len(cart.Items) > 0 {
		return cart, nil
	}

	// 2. Cache miss — read from MongoDB.
	cart, err = w.mongo.FindByUser(ctx, userId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Nothing in DB either — return a fresh empty cart.
			return &models.Cart{UserId: userId, Items: []models.CartItem{}}, nil
		}
		return nil, err
	}

	// 3. Backfill Redis.
	if cacheErr := w.redis.Save(ctx, cart); cacheErr != nil {
		log.Printf("cart cache backfill failed for user %s: %v", userId, cacheErr)
	}

	return cart, nil
}

func (w *WriteThroughCartRepository) Save(ctx context.Context, cart *models.Cart) error {
	// 1. Write to MongoDB first — if this fails, we abort.
	if err := w.mongo.Save(ctx, cart); err != nil {
		return err
	}

	// 2. Write to Redis — failure is non-fatal.
	if err := w.redis.Save(ctx, cart); err != nil {
		log.Printf("cart cache write failed for user %s: %v", cart.UserId, err)
	}

	return nil
}

func (w *WriteThroughCartRepository) Clear(ctx context.Context, userId string) error {
	// 1. Delete from MongoDB first — if this fails, we abort.
	if err := w.mongo.Clear(ctx, userId); err != nil {
		return err
	}

	// 2. Invalidate Redis cache — failure is non-fatal.
	if err := w.redis.Clear(ctx, userId); err != nil {
		log.Printf("cart cache clear failed for user %s: %v", userId, err)
	}

	return nil
}
