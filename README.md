# Fiber Wallet Ledger 🧾

A lightning-fast, secure transactional ledger API built with **Go** and the **Fiber** framework. This microservice manages multi-currency virtual balances and ensures absolute data integrity using strict ACID database transactions via **GORM**.

## 🚀 Architectural Highlights
- **High Performance:** Built on `gofiber/fiber`, optimized for rapid request handling and minimal allocation.
- **ACID Compliant:** Financial transfers utilize strict database transactions (`tx.Begin()`, `tx.Commit()`, `tx.Rollback()`). If any part of a transfer fails (e.g., logging the transaction history), the balance deduction is safely rolled back to prevent fund loss.
- **Multi-Currency Support:** Handles distinct logic and decimal precisions for different coin types (e.g., USDT, BNB).
- **ORM Integration:** Uses `GORM` with SQLite (easily swappable to PostgreSQL or MySQL for production environments).

## 🛠 Tech Stack
- **Framework:** Go Fiber (`gofiber/fiber/v2`)
- **Database ORM:** GORM (`gorm.io/gorm`)
- **Driver:** SQLite (local persistence)

## ⚙️ Setup & Installation

1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod init [github.com/card0re/fiber-wallet-ledger](https://github.com/card0re/fiber-wallet-ledger)
   go get [github.com/gofiber/fiber/v2](https://github.com/gofiber/fiber/v2)
   go get gorm.io/gorm
   go get gorm.io/driver/sqlite
   go mod tidy