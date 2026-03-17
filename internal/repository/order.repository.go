package repository

import (
	"context"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	Collection *mongo.Collection
}

func NewOrderRepository(collection *mongo.Collection) *OrderRepository {
	return &OrderRepository{Collection: collection}
}

func (r *OrderRepository) Save(ctx context.Context, order *models.Order) error {
	_, err := r.Collection.InsertOne(ctx, order)
	return err
}
