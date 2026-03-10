# VAL Jewelry Inventory — Setup Guide

## Prerequisites

- Go 1.22+
- Node.js 18+
- PostgreSQL 14+ (or Docker)

---

## Quick Start (with Docker)

```bash
# 1. Start PostgreSQL (creates DB + runs migrations automatically)
docker-compose up postgres -d

# 2. Start the API server
cp .env.example .env
go run ./cmd/server

# 3. Start the frontend dev server
cd frontend
npm install
npm run dev
# → Open http://localhost:3000
```

---

## Manual Setup

### 1. PostgreSQL

```bash
createdb val_inventory
psql val_inventory < migrations/001_initial.sql
```

### 2. API Server

```bash
cp .env.example .env
# Edit .env with your DATABASE_URL
go run ./cmd/server
# → API runs on http://localhost:8080
```

### 3. Frontend

```bash
cd frontend
npm install
npm run dev          # dev mode → http://localhost:3000
npm run build        # production build → cmd/server/static/
```

---

## Migrate Legacy DBF Data

```bash
# Point at the original VAL directory
go run ./cmd/migrate \
  --dbf /path/to/VAL \
  --dsn "postgres://postgres:postgres@localhost:5432/val_inventory?sslmode=disable"
```

This imports (in order):
1. Account Groups (GRPMST)
2. Account Master (FAMST)
3. Items (ITEM)
4. Lots (LOT)
5. Gold Rates (RATE)
6. Sales (SALE)
7. Purchases (PURCH)

---

## Network / Multi-user

The API server binds to `0.0.0.0:8080` — accessible from any machine on the network.

For the web UI, all users point their browser at:
```
http://<server-ip>:8080
```

For the **desktop app** (Wails), install Wails v2 and run:
```bash
# Install Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Build desktop app
wails build
```

---

## API Reference

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/items | List items (paginated, searchable) |
| POST | /api/v1/items | Create item |
| GET | /api/v1/items/:itcd | Get item by code |
| PUT | /api/v1/items/:itcd | Update item |
| DELETE | /api/v1/items/:itcd | Delete item |
| GET | /api/v1/issues | List issue vouchers |
| POST | /api/v1/issues | Create issue with detail lines |
| GET | /api/v1/receipts | List receipt vouchers |
| POST | /api/v1/receipts | Create receipt with detail lines |
| GET | /api/v1/sales | List sales |
| POST | /api/v1/sales | Create sale |
| GET | /api/v1/purchases | List purchases |
| POST | /api/v1/purchases | Create purchase |
| GET | /api/v1/accounts | List account master |
| GET | /api/v1/accounts/:acCode/address | Get party address |
| GET | /api/v1/lots | List lots |
| GET | /api/v1/rates | List gold rates |
| GET | /api/v1/rates/latest | Latest gold rate |

All list endpoints support `?page=1&per_page=50` pagination.

---

## Project Structure

```
inventory-management/
├── cmd/
│   ├── server/main.go        ← API server entrypoint
│   └── migrate/main.go       ← DBF → PostgreSQL importer
├── internal/
│   ├── db/db.go              ← DB connection
│   ├── models/models.go      ← All domain structs
│   ├── repository/           ← Data access (one file per domain)
│   └── api/                  ← HTTP handlers + router
├── migrations/
│   └── 001_initial.sql       ← Full PostgreSQL schema
├── frontend/
│   ├── src/
│   │   ├── api/client.ts     ← Typed API client
│   │   ├── pages/            ← One page per module
│   │   └── components/       ← Shared components
│   └── package.json
├── docker-compose.yml
└── Dockerfile
```
