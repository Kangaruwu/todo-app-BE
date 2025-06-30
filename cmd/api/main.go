package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()

	app.Get("/", hello);
	app.Use(middleware);

	app.Listen(":8080")
}

func hello(c *fiber.Ctx) error {
	log.Println("Hello handler called")
	return c.SendString("hello wfergerg")
}

func middleware(c *fiber.Ctx) error {
	return c.SendStatus(404)
}