// @title           Food Ordering API
// @version         1.0
// @description     REST API for food ordering with JWT authentication, cart management, and order placement.
// @host            localhost:8000
// @BasePath        /api
// @schemes         http https

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Enter: Bearer <your-jwt-token>

//go:generate swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
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

	_ "food-ordering/docs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {

	config.LoadEnv()
	db := config.LoadMongoConnection()
	redisClient := config.LoadRedisConnection()

	jwtSecret := config.GetEnv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env var is required")
	}

	database.EnsureIndexes(db)

	log.Println("Loading promo code validator (this may take a moment)...")
	promoValidator, err := promo.Load(promoSourcesFromEnv())
	if err != nil {
		log.Fatalf("Failed to load promo validator: %v", err)
	}

	ServerPort := config.GetEnv("SERVER_PORT")
	app := fiber.New()

	// Serve generated swagger.json
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})

	// Serve swagger UI (loads swagger-ui from CDN, points at /swagger/doc.json)
	app.Get("/swagger", func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(swaggerUI)
	})

	// Serve static assets from docs/ (swagger.json reachable directly too)
	app.Use("/docs", static.New("./docs"))

	UserCollection := db.Collection("user")
	ProductCollection := db.Collection("product")
	CartCollection := db.Collection("cart")
	OrderCollection := db.Collection("order")

	authHandler, authService := handlers.InitAuthModule(UserCollection, jwtSecret)
	productHandler := handlers.InitProductModule(ProductCollection)
	cartHandler := handlers.InitCartModule(CartCollection, redisClient)
	orderHandler := handlers.InitOrderModule(OrderCollection, ProductCollection, promoValidator, redisClient)

	appRouter := app.Group("/api")
	routes.AuthRoutes(appRouter, authHandler)
	routes.ProductRoutes(appRouter, productHandler)
	routes.CartRoutes(appRouter, cartHandler, authService)
	routes.OrderRoutes(appRouter, orderHandler, authService, redisClient)

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

	return []string{}
	s1 := config.GetEnv("COUPON_FILE_1")
	s2 := config.GetEnv("COUPON_FILE_2")
	s3 := config.GetEnv("COUPON_FILE_3")
	if s1 == "" && s2 == "" && s3 == "" {
		return nil
	}
	return []string{s1, s2, s3}
}

const swaggerUI = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Food Ordering API — Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({
      url: "/swagger/doc.json",
      dom_id: "#swagger-ui",
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
      layout: "BaseLayout",
      deepLinking: true,
      defaultModelsExpandDepth: 1,
    });
  </script>
</body>
</html>`
