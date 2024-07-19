package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/config"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"github.com/ilhamosaurus/fiber-commerce/util"
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

// @Summary	Register new User
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		user	body		models.RegisterValidation	true	"User"
// @Success	201		{object} handler.Register.RegisterResponse	"User Created"
// @Failure	400		{object}	string						"Invalid fields"
// @Failure	409		{object}	string						"User already exists"
// @Failure	500		{object}	string						"Internal server error"
// @Router		/api/auth/register [post]
func Register(c *fiber.Ctx) error {
	validate := validator.New()
	db := database.DB

	user := &models.RegisterValidation{}
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields"})
	}

	// custom validation for role
	validate.RegisterValidation("role", func(fl validator.FieldLevel) bool {
		// role must containt either "CLIENT" or "MERCHANT"
		return fl.Field().String() == "CLIENT" || fl.Field().String() == "MERCHANT"
	})

	if err := validate.Struct(user); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
	}

	hash, err := util.HashedPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password", "data": err})
	}

	user.Password = hash
	role := models.Role(user.Role)
	fmt.Println(role)
	if err := db.Create(&models.User{Username: user.Username, Password: user.Password, Role: role, Account: &models.Account{Owner: user.Username, Balance: 0}}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(409).JSON(fiber.Map{"error": "Username already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create user", "data": err})
	}

	type RegisterResponse struct {
		Message string `json:"message" example:"User Registered successfully, please login"`
	}

	return c.Status(200).JSON(RegisterResponse{Message: "User Registered successfully, please login"})
}

// @Summary	Login User
// @Tags		Auth
// @Accept		json
// @Produce	json
// @Param		user	body		models.LoginValidation	true	"User"
// @Success	200		{object}	handler.Login.LoginResponse	"User Logged In"
// @Failure	400		{object}	string					"Invalid fields"
// @Failure	500		{object}	string					"Internal server error"
// @Router		/api/auth/login [post]
func Login(c *fiber.Ctx) error {
	validate := validator.New()
	type UserData struct {
		ID       uint        `json:"id"`
		Username string      `json:"username"`
		Role     models.Role `json:"role"`
	}

	input := &models.LoginValidation{}
	var userData UserData

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields"})
	}

	if err := validate.Struct(input); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
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
		userData.Role = userModel.Role
	}

	if !util.CheckPasswordHash(input.Password, userModel.Password) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credetials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = userData.ID
	claims["username"] = userData.Username
	claims["role"] = userData.Role
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to login", "data": err})
	}

	type LoginResponse struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"`
	}

	return c.JSON(LoginResponse{Token: t})
}
