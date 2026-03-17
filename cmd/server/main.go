package main

import (
	"context"
	"fmt"
	"food-ordering/config"
	"food-ordering/internal/api/handlers"
	"food-ordering/internal/api/routes"
	"food-ordering/internal/database"
	"food-ordering/internal/promo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {

	config.LoadEnv()
	db := config.LoadMongoConnection()
	redisClient := config.LoadRedisConnection()

	jwtSecret := config.GetEnv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env var is required")
	}

	// Ensure all MongoDB indexes exist before accepting traffic.
	database.EnsureIndexes(db)

	// Load promo validator at startup. Sources can be overridden via env vars;
	// defaults to downloading the three S3-hosted gz files.
	log.Println("Loading promo code validator (this may take a moment)...")
	promoValidator, err := promo.Load(promoSourcesFromEnv())
	if err != nil {
		log.Fatalf("Failed to load promo validator: %v", err)
	}

	ServerPort := config.GetEnv("SERVER_PORT")
	app := fiber.New()

	UserCollection := db.Collection("user")
	ProductCollection := db.Collection("product")
	CartCollection := db.Collection("cart")
	OrderCollection := db.Collection("order")

	authHandler, authService := handlers.InitAuthModule(UserCollection, jwtSecret)
	productHandler := handlers.InitProductModule(ProductCollection)
	cartHandler := handlers.InitCartModule(CartCollection, redisClient) // write-through: Redis + MongoDB
	orderHandler := handlers.InitOrderModule(OrderCollection, ProductCollection, promoValidator)

	appRouter := app.Group("/api")
	routes.AuthRoutes(appRouter, authHandler)
	routes.ProductRoutes(appRouter, productHandler)
	routes.CartRoutes(appRouter, cartHandler, authService)
	routes.OrderRoutes(appRouter, orderHandler, authService)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Println("shutdown error:", err)
		}
		config.CloseRedisConnection()
		config.CloseMongoConnection()
		log.Println("Server exited properly")
	}()

	if err := app.Listen(fmt.Sprintf(":%s", ServerPort), fiber.ListenConfig{
		EnablePrefork:     false,
		EnablePrintRoutes: true,
	}); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func promoSourcesFromEnv() []string {
	s1 := config.GetEnv("COUPON_FILE_1")
	s2 := config.GetEnv("COUPON_FILE_2")
	s3 := config.GetEnv("COUPON_FILE_3")
	if s1 == "" && s2 == "" && s3 == "" {
		return nil
	}
	return []string{s1, s2, s3}
}
