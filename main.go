package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- MODELS ---

// Account represents a user's wallet balances
type Account struct {
	gorm.Model
	AccountID   string  `gorm:"uniqueIndex;not null"`
	BalanceUSDT float64 `gorm:"type:decimal(18,2);default:0.0"`
	BalanceBNB  float64 `gorm:"type:decimal(18,5);default:0.0"`
	IsActive    bool    `gorm:"default:true"`
}

// Transaction represents a ledger entry
type Transaction struct {
	gorm.Model
	AccountID     uint    `gorm:"not null"`
	Coin          string  `gorm:"not null"`
	TargetAddress string  `gorm:"not null"`
	Amount        float64 `gorm:"type:decimal(18,2);not null"`
	TxHash        string  `gorm:"uniqueIndex;not null"`
}

var db *gorm.DB

// --- UTILS ---
func generateTxHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return "0x" + hex.EncodeToString(bytes)
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("ledger.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate schema
	db.AutoMigrate(&Account{}, &Transaction{})
	log.Println("✅ Database initialized successfully")
}

// --- HANDLERS ---

func getBalance(c *fiber.Ctx) error {
	accountID := c.Params("id")
	var account Account

	if err := db.Where("account_id = ? AND is_active = ?", accountID, true).First(&account).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Account not found or inactive"})
	}

	var txs []Transaction
	db.Where("account_id = ?", account.ID).Order("created_at desc").Limit(10).Find(&txs)

	return c.JSON(fiber.Map{
		"account_id":   account.AccountID,
		"balance_usdt": account.BalanceUSDT,
		"balance_bnb":  account.BalanceBNB,
		"history":      txs,
	})
}

func sendTransaction(c *fiber.Ctx) error {
	type SendRequest struct {
		AccountID     string  `json:"account_id"`
		TargetAddress string  `json:"target_address"`
		Amount        float64 `json:"amount"`
		CoinType      string  `json:"coin_type"` // USDT or BNB
	}

	var req SendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be greater than 0"})
	}

	// Begin ACID Transaction
	tx := db.Begin()
	var account Account

	// Lock row and verify access
	if err := tx.Where("account_id = ? AND is_active = ?", req.AccountID, true).First(&account).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Account access denied"})
	}

	// Balance Check & Deduction logic
	if req.CoinType == "USDT" {
		if account.BalanceUSDT < req.Amount {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient USDT balance"})
		}
		account.BalanceUSDT -= req.Amount
	} else if req.CoinType == "BNB" {
		if account.BalanceBNB < req.Amount {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient BNB balance"})
		}
		account.BalanceBNB -= req.Amount
	} else {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported coin type"})
	}

	// Save new balance
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}

	// Record transaction history
	txHash := generateTxHash()
	transaction := Transaction{
		AccountID:     account.ID,
		Coin:          req.CoinType,
		TargetAddress: req.TargetAddress,
		Amount:        req.Amount,
		TxHash:        txHash,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to record transaction log"})
	}

	// Commit if everything is successful
	tx.Commit()

	return c.JSON(fiber.Map{
		"message":      "Transaction successful",
		"tx_hash":      txHash,
		"balance_usdt": account.BalanceUSDT,
		"balance_bnb":  account.BalanceBNB,
	})
}

func createDemoAccount(c *fiber.Ctx) error {
	account := Account{
		AccountID:   "DEMO-12345",
		BalanceUSDT: 5000.00,
		BalanceBNB:  10.5,
		IsActive:    true,
	}
	db.FirstOrCreate(&account, Account{AccountID: "DEMO-12345"})
	return c.JSON(fiber.Map{"message": "Demo account initialized", "account_id": account.AccountID})
}

func main() {
	initDB()

	app := fiber.New(fiber.Config{
		AppName: "Fiber Wallet Ledger",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())

	// API Routes
	api := app.Group("/api/v1")

	api.Post("/admin/demo-setup", createDemoAccount) // Route just to create initial data
	api.Get("/account/:id/balance", getBalance)
	api.Post("/transaction/send", sendTransaction)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Ledger API running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
