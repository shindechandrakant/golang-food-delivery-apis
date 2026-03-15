package main

import (
	"context"
	"fmt"
	"food-ordering/config"

	//"food-ordering/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Image struct {
	Thumbnail string `bson:"thumbnail"`
	Mobile    string `bson:"mobile"`
	Tablet    string `bson:"tablet"`
	Desktop   string `bson:"desktop"`
}

type Product struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Cuisines    []string      `bson:"cuisines,omitempty"`
	Category    string        `bson:"category"`
	Price       float64       `bson:"price"`
	Description string        `bson:"description,omitempty"`
	Rating      float32       `bson:"rating"`
	Image       []Image       `bson:"image"`
}

func main() {

	config.LoadEnv()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoDbConnection := config.LoadMongoConnection()
	ProductCollection := mongoDbConnection.Collection("product")

	products := []interface{}{
		Product{
			Name:        "Margherita Pizza",
			Cuisines:    []string{"Italian"},
			Category:    "Pizza",
			Price:       8.99,
			Description: "Classic pizza with tomato sauce and mozzarella",
			Rating:      4.5,
			Image: []Image{
				{
					Thumbnail: "margherita-thumb.jpg",
					Mobile:    "margherita-mobile.jpg",
					Tablet:    "margherita-tablet.jpg",
					Desktop:   "margherita-desktop.jpg",
				},
			},
		},
		Product{
			Name:        "Chicken Biryani",
			Cuisines:    []string{"Indian", "Hyderabadi"},
			Category:    "Rice",
			Price:       12.5,
			Description: "Spicy basmati rice cooked with chicken",
			Rating:      4.7,
			Image: []Image{
				{
					Thumbnail: "biryani-thumb.jpg",
					Mobile:    "biryani-mobile.jpg",
					Tablet:    "biryani-tablet.jpg",
					Desktop:   "biryani-desktop.jpg",
				},
			},
		},
		Product{
			Name:        "Veg Hakka Noodles",
			Cuisines:    []string{"Chinese"},
			Category:    "Noodles",
			Price:       7.25,
			Description: "Stir fried noodles with vegetables",
			Rating:      4.2,
			Image: []Image{
				{
					Thumbnail: "noodles-thumb.jpg",
					Mobile:    "noodles-mobile.jpg",
					Tablet:    "noodles-tablet.jpg",
					Desktop:   "noodles-desktop.jpg",
				},
			},
		},
		Product{
			Name:        "Paneer Butter Masala",
			Cuisines:    []string{"Indian", "Punjabi"},
			Category:    "Curry",
			Price:       10.75,
			Description: "Creamy tomato curry with paneer",
			Rating:      4.6,
			Image: []Image{
				{
					Thumbnail: "paneer-thumb.jpg",
					Mobile:    "paneer-mobile.jpg",
					Tablet:    "paneer-tablet.jpg",
					Desktop:   "paneer-desktop.jpg",
				},
			},
		},
		Product{
			Name:        "Chocolate Brownie",
			Cuisines:    []string{"Dessert"},
			Category:    "Dessert",
			Price:       5.5,
			Description: "Rich chocolate brownie served warm",
			Rating:      4.8,
			Image: []Image{
				{
					Thumbnail: "brownie-thumb.jpg",
					Mobile:    "brownie-mobile.jpg",
					Tablet:    "brownie-tablet.jpg",
					Desktop:   "brownie-desktop.jpg",
				},
			},
		},
	}

	result, err := ProductCollection.InsertMany(ctx, products)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted documents:", result.InsertedIDs)
}
