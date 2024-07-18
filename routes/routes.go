package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/handler"
	"github.com/ilhamosaurus/fiber-commerce/middleware"
)

func SetupRoutes(app *fiber.App) {
	// api global set prefix
	api := app.Group("/api")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// product routes
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:code", handler.GetProduct)

	// auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// transaction routes
	transaction := api.Group("/transaction")
	transaction.Use(middleware.Protected())
	transaction.Get("/balance", handler.GetBalance)
	transaction.Post("/topup", handler.Topup)
}
