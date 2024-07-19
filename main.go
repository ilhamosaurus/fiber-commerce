package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/ilhamosaurus/fiber-commerce/database"
	_ "github.com/ilhamosaurus/fiber-commerce/docs"
	"github.com/ilhamosaurus/fiber-commerce/routes"
)

// @title			Fiber-Mini Commerce
// @version		1.0
// @description	API Documentation for Fiber-Mini Commerce which is an e-commerce application.
// @description	Where user either can be a client or merchant. Client can buy product and merchant can sell product.
// @description	PS: Authorization cannot be used in this project because OpenAPi 2.0 does not support Bearer Token.
// @host			localhost:3000
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @BasePath		/
func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New())
	database.ConnectDb()

	app.Get("/api-docs/*", swagger.HandlerDefault)

	routes.SetupRoutes(app)

	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
