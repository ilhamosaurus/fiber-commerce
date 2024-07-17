package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/routes"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())
	database.ConnectDb()

	routes.SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
