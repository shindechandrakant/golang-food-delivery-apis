package main

import (
	"context"
	"fmt"
	"food-ordering/config"
	"food-ordering/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     string
}

var users = []seedUser{
	{
		Name:     "Admin User",
		Email:    "admin@foodapp.com",
		Password: "admin123",
		Role:     models.RoleAdmin,
	},
	{
		Name:     "Alice Johnson",
		Email:    "alice@example.com",
		Password: "password123",
		Role:     models.RoleUser,
	},
	{
		Name:     "Bob Smith",
		Email:    "bob@example.com",
		Password: "password123",
		Role:     models.RoleUser,
	},
	{
		Name:     "Carol White",
		Email:    "carol@example.com",
		Password: "password123",
		Role:     models.RoleUser,
	},
}

func main() {
	config.LoadEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db := config.LoadMongoConnection()
	col := db.Collection("user")

	inserted, skipped := 0, 0

	for _, u := range users {
		// Skip if email already exists.
		var existing models.User
		err := col.FindOne(ctx, bson.M{"email": u.Email}).Decode(&existing)
		if err == nil {
			fmt.Printf("  skip  %s (%s) — already exists\n", u.Email, u.Role)
			skipped++
			continue
		}
		if err != mongo.ErrNoDocuments {
			log.Fatalf("error checking %s: %v", u.Email, err)
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("bcrypt error for %s: %v", u.Email, err)
		}

		doc := models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashed),
			Role:     u.Role,
		}

		if _, err := col.InsertOne(ctx, doc); err != nil {
			log.Fatalf("insert error for %s: %v", u.Email, err)
		}

		fmt.Printf("  added %s (%s)\n", u.Email, u.Role)
		inserted++
	}

	fmt.Printf("\nDone — inserted: %d, skipped: %d\n", inserted, skipped)
}
