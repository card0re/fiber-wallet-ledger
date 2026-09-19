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

## ⚙️ Run

```bash
git clone https://github.com/card0re/fiber-wallet-ledger.git
cd fiber-wallet-ledger
go run .            # listens on :3000 (override with PORT), creates ledger.db
```

## 📡 API

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/admin/demo-setup` | Create demo account `DEMO-12345` |
| GET | `/api/v1/account/:id/balance` | Get balances |
| POST | `/api/v1/transaction/send` | Transfer funds (atomic) |

```bash
curl -X POST http://localhost:3000/api/v1/admin/demo-setup
curl -X POST http://localhost:3000/api/v1/transaction/send \
  -H "Content-Type: application/json" \
  -d '{"account_id":"DEMO-12345","target_address":"0xabc","amount":10,"coin_type":"USDT"}'
curl http://localhost:3000/api/v1/account/DEMO-12345/balance
```
