package repository

import (
	"context"
	"food-ordering/internal/dto"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductRepository struct {
	Collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) *ProductRepository {
	return &ProductRepository{Collection: collection}
}

// FindWithFilter fetches products applying optional name search, category,
// price range, and sorting. Zero-value fields in the filter are ignored.
func (r *ProductRepository) FindWithFilter(ctx context.Context, f dto.ProductFilter) ([]models.Product, error) {
	query := bson.D{}

	if f.Name != "" {
		query = append(query, bson.E{Key: "name", Value: bson.M{"$regex": f.Name, "$options": "i"}})
	}

	if f.Category != "" {
		query = append(query, bson.E{Key: "category", Value: f.Category})
	}

	if f.MinPrice > 0 || f.MaxPrice > 0 {
		priceFilter := bson.M{}
		if f.MinPrice > 0 {
			priceFilter["$gte"] = f.MinPrice
		}
		if f.MaxPrice > 0 {
			priceFilter["$lte"] = f.MaxPrice
		}
		query = append(query, bson.E{Key: "price", Value: priceFilter})
	}

	opts := options.Find().SetSort(buildSort(f.SortBy, f.SortOrder))

	cursor, err := r.Collection.Find(ctx, query, opts)
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

// buildSort returns a bson.D sort document.
// Allowed sortBy values: "name", "price", "rating". Defaults to "name".
// Allowed sortOrder values: "asc", "desc". Defaults to "asc".
func buildSort(sortBy, sortOrder string) bson.D {
	field := "name"
	switch sortBy {
	case "price", "rating":
		field = sortBy
	}

	order := 1
	if sortOrder == "desc" {
		order = -1
	}

	return bson.D{{Key: field, Value: order}}
}
