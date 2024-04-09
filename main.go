package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"search_engine/db"
	"search_engine/routes"
	"search_engine/utils"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("cannot find environment variables")
	}
	log.Println("Env loaded")

	db.InitDb()
	log.Println("Database Initialized")
}

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	panic("cannot find environment variables")
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	app.Use(compress.New())

	routes.SetRoutes(app)

	utils.StartCronJobs() // TODO: see if it is possible to put this in the init function

	// Start the server and listen for a shutdown
	go func() {
		if err := app.Listen(port); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // Block the main thread until a signal is received

	_ = app.Shutdown()
	fmt.Println("Shutting down server")
}
