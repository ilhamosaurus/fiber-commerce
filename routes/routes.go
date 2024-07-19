package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/handler"
	"github.com/ilhamosaurus/fiber-commerce/middleware"
)

func SetupRoutes(app *fiber.App) {

	// api global set prefix
	api := app.Group("/api")

	// auth routes
	auth := api.Group("/auth")
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// product routes
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:code", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Put("/:code", middleware.Protected(), handler.UpdateProduct)
	product.Delete("/:code", middleware.Protected(), handler.DeleteProduct)

	// transaction routes
	transaction := api.Group("/transaction")
	transaction.Use(middleware.Protected())
	transaction.Get("/balance", handler.GetBalance)
	transaction.Post("/topup", handler.Topup)
	transaction.Get("/history", handler.GetOrders)
	transaction.Post("/payment", handler.Payment)
}
