package main

import (
	"context"
	"fmt"
	"food-ordering/config"
	"food-ordering/internal/api/handlers"
	"food-ordering/internal/api/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
)

func main() {

	config.LoadEnv()
	mongoDbConnection := config.LoadMongoConnection()
	_ = config.LoadRedisConnection()

	ServerPort := config.GetEnv("SERVER_PORT")
	app := fiber.New()

	ProductCollection := mongoDbConnection.Collection("product")
	productHandler := handlers.InitProductModule(ProductCollection)
	appRouter := app.Group("/api")
	routes.ProductRoutes(appRouter, productHandler)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := app.ShutdownWithContext(ctx)
		if err != nil {
			log.Println("unable to close the app", err)
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
