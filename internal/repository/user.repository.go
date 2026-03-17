package repository

import (
	"context"
	"food-ordering/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{Collection: collection}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Save(ctx context.Context, user *models.User) error {
	result, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	user.Id = result.InsertedID.(bson.ObjectID)
	return nil
}
