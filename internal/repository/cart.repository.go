package repository

import (
	"context"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type CartRepository struct {
	Collection *mongo.Collection
}

func NewCartRepository(collection *mongo.Collection) *CartRepository {
	return &CartRepository{
		Collection: collection,
	}
}

func (c *CartRepository) FindByUser(ctx context.Context, userId string) (*models.Cart, error) {
	var cart models.Cart
	err := c.Collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (c *CartRepository) Save(ctx context.Context, cart *models.Cart) error {

	opts := options.UpdateOne().SetUpsert(true)
	_, err := c.Collection.UpdateOne(ctx,
		bson.M{"userId": cart.UserId},
		bson.M{"$set": cart},
		opts)

	return err
}

func (c *CartRepository) Clear(ctx context.Context, userId string) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"userId": userId})
	return err
}
