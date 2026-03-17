package repository

import (
	"context"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// CartRepository is the persistence contract for the cart.
type CartRepository interface {
	FindByUser(ctx context.Context, userId string) (*models.Cart, error)
	Save(ctx context.Context, cart *models.Cart) error
	Clear(ctx context.Context, userId string) error
}

// MongoCartRepository is a MongoDB-backed CartRepository (kept as fallback).
type MongoCartRepository struct {
	Collection *mongo.Collection
}

func NewMongoCartRepository(collection *mongo.Collection) *MongoCartRepository {
	return &MongoCartRepository{Collection: collection}
}

func (c *MongoCartRepository) FindByUser(ctx context.Context, userId string) (*models.Cart, error) {
	var cart models.Cart
	err := c.Collection.FindOne(ctx, bson.M{"userId": userId}).Decode(&cart)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (c *MongoCartRepository) Save(ctx context.Context, cart *models.Cart) error {
	opts := options.UpdateOne().SetUpsert(true)
	_, err := c.Collection.UpdateOne(ctx,
		bson.M{"userId": cart.UserId},
		bson.M{"$set": cart},
		opts)
	return err
}

func (c *MongoCartRepository) Clear(ctx context.Context, userId string) error {
	_, err := c.Collection.DeleteOne(ctx, bson.M{"userId": userId})
	return err
}
