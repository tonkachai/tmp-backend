package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"tmp-backend/db"
	"tmp-backend/handlers"
	"tmp-backend/utils"
)

func main() {
	db.Init()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Hello World"})
	})

	api := app.Group("/")
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	api.Get("/me", utils.JWTMiddleware, handlers.Me)

	app.Get("/swagger.json", func(c *fiber.Ctx) error {
		return c.SendFile("./swagger.json")
	})

	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.SendFile("./swagger-ui.html")
	})

	log.Fatal(app.Listen(":3000"))
}
