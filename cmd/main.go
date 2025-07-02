package main

import (
	"context"
	"go-backend-todo/internal/db"

	"go-backend-todo/internal/api/handlers"
	"go-backend-todo/internal/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", handlers.Hello)
	app.Use(middlewares.NotFound)

	db, err := db.Connect()
	if err != nil {
		panic(err)
	}
	defer db.Close(context.Background())

	app.Listen(":8080")

}


