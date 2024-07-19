package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/models"
)

type CurUser struct {
	ID        uint
	Username string
	Role models.Role
}

func CurrentUser(c *fiber.Ctx) CurUser {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	role := claims["role"].(string)
	return CurUser{
		ID:        uint(claims["id"].(float64)),
		Username:  username,
		Role:      models.Role(role),
	}
}