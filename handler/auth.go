package handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/config"
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

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
	}

	input := new(LoginInput)
	var userData UserData

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields"})
	}

	userModel, err := new(models.User), *new(error)

	userModel, err = getUserByUsername(input.Username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to find user", "data": err})
	}

	if userModel == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid credetials", "data": err})
	} else {
		userData.ID = userModel.ID
		userData.Username = userModel.Username
	}

	if !CheckPasswordHash(input.Password, userModel.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credetials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userData.ID
	claims["username"] = userData.Username
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to login", "data": err})
	}

	return c.JSON(fiber.Map{"token": t})
}
