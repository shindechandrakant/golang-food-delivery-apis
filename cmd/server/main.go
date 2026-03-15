package main

import (
	"fmt"
	"food-ordering/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
)

func main() {

	config.LoadEnv()
	_ = config.LoadMongoConnection()
	_ = config.LoadRedisConnection()

	SERVER_PORT := config.GetEnv("SERVER_PORT")
	app := fiber.New()

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", SERVER_PORT), fiber.ListenConfig{
			EnablePrefork:     true,
			EnablePrintRoutes: true,
		}); err != nil {
			log.Fatal()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	err := app.Shutdown()
	if err != nil {
		log.Println("unable to close the app", err)
	}
	config.CloseRedisConnection()
	config.CloseMongoConnection()
	log.Println("Server exited properly")
}
