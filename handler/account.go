package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/gorm"
)

type AccountData struct {
	ID      uint    `json:"id"`
	Owner   string  `json:"owner"`
	Balance float64 `json:"balance"`
}

type TransactionUtil struct {
	Invoice string       `json:"invoice"`
	Account *AccountData `json:"account"`
}

func GetAccountByUsername(username string) (*AccountData, error) {
	db := database.DB

	var account models.Account
	if err := db.Where(&models.Account{Owner: username}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	accountData := AccountData{
		ID:      account.ID,
		Owner:   account.Owner,
		Balance: account.Balance,
	}

	return &accountData, nil
}

func GetInvNumber(username string) (*TransactionUtil, error) {
	db := database.DB

	account, err := GetAccountByUsername(username)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}
	var count int64
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())
	if err := db.Model(&models.Order{}).Where("account_id = ? AND created_at >= ?", account.ID, today).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		count = 0
	}

	invNumber := fmt.Sprintf("INV%s-%04d", today.Format("02012006"), count+1)

	transactionUtil := TransactionUtil{
		Invoice: invNumber,
		Account: account,
	}

	return &transactionUtil, nil
}

// @Summary Get Account Balance
// @Tags Transaction
// @Description Get Account Balance
// @Security Bearer
// @Produce json
// @Success 200 {object} handler.GetBalance.BalanceResponse
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Account not found"
// @Failure 500 {object} string "Failed to get balance"
// @Router /api/transaction/balance [get]
func GetBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)

	account, err := GetAccountByUsername(username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get balance", "data": err})
	}

	if account == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found", "data": nil})
	}

	type BalanceResponse struct {
		Owner   string  `json:"owner"`
		Balance float64 `json:"balance"`
	}

	return c.Status(200).JSON(BalanceResponse{account.Owner, account.Balance})
}
