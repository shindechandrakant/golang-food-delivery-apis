package repository

import (
	"context"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	Collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{Collection: collection}
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindById(ctx context.Context, id bson.ObjectID) (*models.Product, error) {

	var product models.Product
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)

	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByIds(ctx context.Context, ids []bson.ObjectID) ([]models.Product, error) {

	cursor, err := r.Collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}
