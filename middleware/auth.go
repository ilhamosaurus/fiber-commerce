package middleware

import (
	jwtWare "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/config"
)

// protected routes
func Protected() fiber.Handler {
	return jwtWare.New(jwtWare.Config{
		SigningKey:   jwtWare.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}
	return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired JWT"})
}
