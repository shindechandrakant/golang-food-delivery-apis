package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// EnsureIndexes creates all necessary indexes for the database.
// It is idempotent — safe to call on every startup.
func EnsureIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ensureUserIndexes(ctx, db.Collection("user"))
	ensureProductIndexes(ctx, db.Collection("product"))
	ensureOrderIndexes(ctx, db.Collection("order"))
}

func ensureUserIndexes(ctx context.Context, col *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("email_unique"),
		},
	}
	createIndexes(ctx, col, indexes)
}

func ensureProductIndexes(ctx context.Context, col *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "category", Value: 1}},
			Options: options.Index().SetName("category_asc"),
		},
		{
			Keys:    bson.D{{Key: "price", Value: 1}},
			Options: options.Index().SetName("price_asc"),
		},
		{
			// Text index for name/description search.
			Keys:    bson.D{{Key: "name", Value: "text"}, {Key: "description", Value: "text"}},
			Options: options.Index().SetName("product_text_search"),
		},
	}
	createIndexes(ctx, col, indexes)
}

func ensureOrderIndexes(ctx context.Context, col *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("uuid_unique"),
		},
	}
	createIndexes(ctx, col, indexes)
}

func createIndexes(ctx context.Context, col *mongo.Collection, indexes []mongo.IndexModel) {
	if len(indexes) == 0 {
		return
	}
	names, err := col.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		log.Printf("index creation warning on %s: %v", col.Name(), err)
		return
	}
	log.Printf("indexes ready on %s: %v", col.Name(), names)
}
