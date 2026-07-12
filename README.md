# MunchiesGo

A multi-vendor food delivery backend written in Go — vendors, brands, menus, orders and riders behind a single, type-safe REST API.

## Overview

A food delivery marketplace is more than "a restaurant with a menu". One vendor company can operate several brands, each brand has its own menu, opening hours and delivery zones, and a single customer basket may span multiple brands and must be split into per-brand tickets.

MunchiesGo models that domain properly. It is a production-shaped REST API built on Gin and PostgreSQL, with SQL-first, compile-time-checked data access (sqlc), PASETO-based authentication carrying role and vendor scope, versioned migrations, and transactional order placement.

## Features

- **Vendor management** — create, list, search, update and delete vendor companies, with payment info and social links stored as `jsonb`.
- **Staff accounts & authentication** — username/password login issuing a PASETO access token; role-based staff model (`admin`, `manager`, `reporting`, `operations`, `finance`, `support`); passwords hashed with bcrypt; bearer-token middleware protecting every management route.
- **Brand management** — multiple brands per vendor, each with slug, logo, banner, geo `Point` location, address, availability flag, rating and status.
- **Operating hours** — per-brand weekly schedule (day of week plus start/end time) with full CRUD.
- **Menu structure** — global catalogue categories plus per-brand menu sections, and products with image, price, availability and generated slugs.
- **Product variants** — nested option groups and variant items (size, add-ons) with `minSelect` / `maxSelect` rules, created in a single bulk upload path.
- **Inventory** — optional per-product stock tracking with purchase price, selling price and expiry date, backed by an inventory history ledger (`input`, `sold`, `restock`, `damage`).
- **Customers** — signup and saved delivery addresses with geo coordinates, label, apartment, floor and area.
- **Order placement** — one database transaction that creates the order, resolves each item's real product, brand and price server-side (never trusting client-supplied prices), writes the order lines, and fans a multi-brand basket out into per-brand `brand_orders`.
- **Order lifecycle & audit trail** — status machine (`PENDING → ACCEPTED / REJECTED / CANCELLED → READY → PICKED → DELIVERED`), payment method (`CASH`, `ONLINE`) and payment status, plus an `order_logs` table recording actor, action and before/after values.
- **Delivery** — rider accounts, rider-to-order assignment, address snapshots frozen at order time, and delivery zones defined as PostgreSQL `Polygon` boundaries mapped to brands.

## Tech Stack

| Concern | Choice |
| --- | --- |
| Language | Go 1.18 |
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| Database | PostgreSQL 14 (with `Point` / `Polygon` geometry columns) |
| Data access | [sqlc](https://sqlc.dev) — type-safe Go generated from raw SQL (no ORM) |
| Driver | `lib/pq` |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Authentication | [PASETO](https://github.com/o1egl/paseto) tokens |
| Password hashing | `golang.org/x/crypto/bcrypt` |
| Configuration | [Viper](https://github.com/spf13/viper) — `app.env` + environment variables |
| Validation | `go-playground/validator` via Gin binding tags |
| Testing | `testify` |
| Tooling | Docker (PostgreSQL), Make |

## Architecture

A clean layered design: a hand-written business layer wrapped around generated, type-safe data access.

- **`api/`** — HTTP layer. Gin router and route groups, auth middleware, request binding and validation, one handler file per resource, plus cross-cutting helpers (slug generation, geo `Point` marshalling, brand prefixes).
- **`db/`** — persistence layer. Hand-written SQL in `db/query/` is compiled by sqlc into `db/sqlc/`. A `Store` interface extends the generated `Querier` with hand-written, transaction-backed operations (`PlaceOrder`, `BulkProductUpload`, `FullProduct`) and batch loaders (`ProductsByIds`, `ProductVariantsByProductIds`, …) that avoid N+1 queries.
- **`token/`** — a pluggable `Maker` interface for token creation and verification, implemented with PASETO. The payload carries `userId`, `vendorId` and `role`, so authorization decisions need no extra round trip.
- **`util/`** — configuration loading, password hashing, random helpers.

```
.
├── main.go                      # entrypoint: config → DB → store → server
├── Makefile                     # postgres, migrate, sqlc, test, server targets
├── sqlc.yaml                    # sqlc codegen config
├── api/
│   ├── server.go                # router setup and route groups
│   ├── middleware.go            # bearer-token auth middleware
│   ├── vendor.go
│   ├── vendor_user.go           # staff creation and login
│   ├── brand.go
│   ├── brand_category.go
│   ├── brand_operation_times.go
│   ├── brand_zone.go
│   ├── category.go
│   ├── product.go
│   ├── customer.go
│   ├── order.go
│   ├── rider.go
│   ├── zone.go
│   └── utils.go                 # Point type, slug and prefix helpers
├── db/
│   ├── migration/               # golang-migrate up/down SQL
│   ├── query/                   # source SQL for sqlc
│   └── sqlc/                    # generated models and queries, Store, PlaceOrder tx
├── token/
│   ├── maker.go                 # Maker interface
│   ├── payload.go
│   └── paseto_maker.go
└── util/
    ├── config.go                # Viper config
    ├── password.go              # bcrypt hash and verify
    └── random.go
```

## Data Model

Core entities, defined in `db/migration/000001_init_schema.up.sql`:

| Entity | Purpose |
| --- | --- |
| `vendors` | Vendor company: name, contact, `payment_info` and `social_links` (jsonb) |
| `vendor_users` | Staff logins scoped to a vendor: hashed password, role |
| `brands` | Storefront owned by a vendor: slug, logo, banner, `location` (Point), address, availability, rating, status |
| `operation_times` | Per-brand weekly opening hours (`day_of_week`, `start_time`, `end_time`) |
| `zones`, `brand_zones` | Delivery zones as `Polygon` boundaries, mapped many-to-many onto brands |
| `categories`, `brand_categories` | Global catalogue categories, and each brand's own menu sections |
| `products` | Menu item: type, price, image, availability, `use_inventory` flag |
| `product_variants`, `product_variant_items` | Option groups (`min_select` / `max_select`) and their selectable items |
| `product_inventory`, `inventory_history` | Stock with purchase/selling price and expiry, plus a ledger of stock movements |
| `customers`, `customer_addresses` | Customer profile and saved addresses with geo points |
| `orders`, `order_items` | Basket header (total, discount, service charge, payable, payment and order status) and its lines |
| `brand_orders` | Per-brand split of a multi-brand order, with its own status and total |
| `delivery_address` | Address snapshot frozen at the moment the order is placed |
| `riders`, `rider_assign` | Rider accounts and their assignment to orders |
| `order_logs` | Audit trail: actor, action type, and previous/current values (jsonb) |

The schema ships with foreign keys throughout and purpose-built indexes, including composite indexes such as `(brand_id, category_id)` and `(brand_id, day_of_week)`.

## Getting Started

### Prerequisites

- Go 1.18 or newer
- Docker (for PostgreSQL)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [sqlc](https://sqlc.dev) CLI (only needed if you change SQL)

### 1. Configure

Create an `app.env` file in the project root (it is git-ignored):

```env
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:admin@localhost:6500/munchies_backend?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012
TOKEN_DURATION=15m
```

| Variable | Description |
| --- | --- |
| `DB_DRIVER` | Database driver name (`postgres`) |
| `DB_SOURCE` | PostgreSQL connection string |
| `SERVER_ADDRESS` | Host and port the HTTP server binds to |
| `TOKEN_SYMMETRIC_KEY` | PASETO symmetric key — must be exactly 32 characters |
| `TOKEN_DURATION` | Access-token lifetime, e.g. `15m` |

Viper reads `app.env` and also honours matching environment variables.

### 2. Database

```bash
make postgres     # run PostgreSQL 14 in Docker on port 6500
make createdb     # create the munchies_backend database
make migrateup    # apply migrations
```

Roll back with `make migratedown`. Scaffold a new migration with `make newmigration name=<name>`.

### 3. Run

```bash
make server       # go run main.go
```

The API is then served at `SERVER_ADDRESS` under the `/api` prefix.

### 4. Regenerate data access (optional)

After editing SQL in `db/query/` or `db/migration/`:

```bash
make sqlc
```

## API Reference

All routes are prefixed with `/api`. Every route except staff login and customer signup requires an `Authorization: Bearer <token>` header.

### Authentication

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/vendor_user/login` | Log in a staff user; returns an access token and the user |

### Vendors

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/vendor/` | Create a vendor |
| `GET` | `/api/vendor/` | List vendors |
| `GET` | `/api/vendor/search` | Search vendors |
| `GET` | `/api/vendor/:id` | Get a vendor |
| `PUT` | `/api/vendor/:id` | Update a vendor |
| `DELETE` | `/api/vendor/:id` | Delete a vendor |

### Vendor Users

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/vendor_user/` | Create a staff user within the caller's vendor |
| `GET` | `/api/vendor_user/` | List staff users |
| `GET` | `/api/vendor_user/:id` | Get a staff user |
| `PUT` | `/api/vendor_user/:id` | Update a staff user |
| `DELETE` | `/api/vendor_user/:id` | Delete a staff user |

### Brands

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/brand/` | Create a brand |
| `GET` | `/api/brand/` | List brands |
| `GET` | `/api/brand/:id` | Get a brand |
| `PUT` | `/api/brand/:id` | Update a brand |
| `DELETE` | `/api/brand/:id` | Delete a brand |

### Brand Operating Hours

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/api/brand/:id/operation_time` | List a brand's operating hours |
| `POST` | `/api/brand/:id/operation_time` | Add an operating-hours slot |
| `PUT` | `/api/brand/operation_time/:id` | Update a slot |
| `DELETE` | `/api/brand/operation_time/:id` | Delete a slot |

### Brand Menu Categories

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/brand/:id/category` | Create a menu category for a brand |
| `GET` | `/api/brand/:id/category` | List a brand's menu categories |
| `GET` | `/api/brand/:id/category/:categoryId` | Get a brand menu category |
| `POST` | `/api/brand/:id/category/:categoryId` | Update a brand menu category |
| `DELETE` | `/api/brand/:id/category/:categoryId` | Delete a brand menu category |

### Categories

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/category/` | Create a global category |
| `GET` | `/api/category/` | List categories |
| `GET` | `/api/category/:id` | Get a category |
| `PUT` | `/api/category/:id` | Update a category |
| `DELETE` | `/api/category/:id` | Delete a category |

### Products

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/product/` | Create a product together with its nested variants and variant items |
| `GET` | `/api/product/:id/details` | Get a product with its variants and items |

### Customers

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/customer/signup` | Register a customer (public) |

### Orders

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/api/order/` | Place an order; prices are resolved server-side and the basket is split per brand |

## Testing

```bash
make test         # go test -v -cover ./...
```

Test suites use `testify` and run against a test server built from the same `Store` interface used in production:

- **`db/sqlc/`** — integration tests for vendor, vendor-user, brand, category and product queries (require a running, migrated database).
- **`token/`** — PASETO token creation, verification, expiry and invalid-token handling.
- **`util/`** — bcrypt password hashing and verification.
- **`api/`** — auth middleware behaviour for valid, missing, malformed and expired tokens.

## Roadmap

Scaffolded and in progress: customer OTP login and profile endpoints, rider APIs and dispatch, delivery-zone geo lookup, discount and service-charge calculation, and admin-scoped authorization on catalogue routes.
