package repository

import (
	"context"
	"encoding/json"
	"errors"
	"food-ordering/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	cartKeyPrefix = "cart:"
	cartTTL       = 30 * 24 * time.Hour // 30 days
)

// RedisCartRepository stores carts in Redis as JSON.
// Key pattern: cart:{userId}
type RedisCartRepository struct {
	client *redis.Client
}

func NewRedisCartRepository(client *redis.Client) *RedisCartRepository {
	return &RedisCartRepository{client: client}
}

func (r *RedisCartRepository) FindByUser(ctx context.Context, userId string) (*models.Cart, error) {
	data, err := r.client.Get(ctx, cartKey(userId)).Bytes()
	if errors.Is(err, redis.Nil) {
		// Return an empty cart — not an error, just nothing stored yet.
		return &models.Cart{UserId: userId, Items: []models.CartItem{}}, nil
	}
	if err != nil {
		return nil, err
	}

	var cart models.Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *RedisCartRepository) Save(ctx context.Context, cart *models.Cart) error {
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, cartKey(cart.UserId), data, cartTTL).Err()
}

func (r *RedisCartRepository) Clear(ctx context.Context, userId string) error {
	return r.client.Del(ctx, cartKey(userId)).Err()
}

func cartKey(userId string) string {
	return cartKeyPrefix + userId
}
