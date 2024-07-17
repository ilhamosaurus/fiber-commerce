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

	// auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// protected routes
	hello := api.Group("/hello")
	hello.Use(middleware.Protected())
	hello.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, from protected route!")
	})
}
