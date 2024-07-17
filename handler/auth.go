package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func getUserByUsername(username string) (*models.User, error) {
	db := database.DB
	var user models.User
	if err := db.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func hashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func Register(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields"})
	}

	hash, err := hashedPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password", "data": err})
	}

	user.Password = hash
	if err := db.Create(&models.User{Username: user.Username, Password: user.Password, Account: &models.Account{Owner: user.Username, Balance: 0}}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user", "data": err})
	}

	return c.Status(200).JSON(fiber.Map{"message": "User Registered successfully, please login"})
}
