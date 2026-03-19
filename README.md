# Food Ordering API

A RESTful food ordering backend built with Go, Fiber, MongoDB, and Redis.

## Tech Stack

- **Runtime:** Go 1.25
- **Framework:** Fiber v3
- **Database:** MongoDB
- **Cache:** Redis (write-through cart caching)
- **Auth:** JWT
- **Docs:** Swagger (swaggo)

---

## Folder Structure

```
├── cmd
│   ├── server
│   │   └── main.go           # application entrypoint
│   ├── seed
│   │   └── main.go           # seed products
│   └── seed-users
│       └── main.go           # seed users
│
├── config
│   ├── env.go                # load env vars
│   ├── mongoDB.go            # MongoDB connection
│   └── redis.go              # Redis connection
│
├── internal
│   ├── api
│   │   ├── routes
│   │   │   ├── auth.routes.go
│   │   │   ├── product.routes.go
│   │   │   ├── cart.routes.go
│   │   │   └── order.routes.go
│   │   │
│   │   └── handlers
│   │       ├── auth.handler.go
│   │       ├── product.handler.go
│   │       ├── cart.handler.go
│   │       └── order.handler.go
│   │
│   ├── services
│   │   ├── auth.service.go
│   │   ├── product.service.go
│   │   ├── cart.service.go
│   │   └── order.service.go
│   │
│   ├── repository
│   │   ├── product.repository.go
│   │   ├── cart.repository.go
│   │   ├── cart_redis.repository.go      # Redis-only cart reads
│   │   ├── cart_writethrough.repository.go  # write-through cache layer
│   │   ├── order.repository.go
│   │   └── user.repository.go
│   │
│   ├── models
│   │   ├── Product.go
│   │   ├── Cart.go
│   │   ├── Order.go
│   │   └── User.go
│   │
│   ├── middleware
│   │   ├── auth.go           # API key auth
│   │   ├── jwt.go            # JWT validation
│   │   ├── idempotency.go    # idempotency key enforcement
│   │   ├── logger.go         # request logging
│   │   └── uuid.go           # request ID injection
│   │
│   ├── dto
│   │   ├── auth.dto.go
│   │   ├── product.dto.go
│   │   ├── cart.dto.go
│   │   ├── order.dto.go
│   │   └── swagger_types.go
│   │
│   ├── database
│   │   └── indexes.go        # MongoDB index setup
│   │
│   ├── promo
│   │   └── validator.go      # coupon/promo validation
│   │
│   └── utils
│       └── response.go       # standard JSON response helpers
│
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── .env
├── go.mod
└── README.md
```

---

## Setup

### Prerequisites

- Go 1.25+
- MongoDB running on `localhost:27017`
- Redis running on `localhost:6379`

### 1. Install dependencies

```bash
go mod download
```

### 2. Configure environment

Create a `.env` file in the project root:

```env
SERVER_PORT=8000

DB_URI=mongodb://localhost:27017
DB_NAME=foodOrdering

REDIS_URL=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=your-secret-here
API_KEY=your-api-key-here

# Paths to coupon bloom filter files (optional)
COUPON_FILE_1=/path/to/couponbase.gz
COUPON_FILE_2=/path/to/couponbase.gz
COUPON_FILE_3=/path/to/couponbase.gz
```

### 3. Run the server

```bash
go run cmd/server/main.go
```

The API will be available at `http://localhost:8000`.

### 4. Seed data (optional)

```bash
# Seed products
go run cmd/seed/main.go

# Seed users
go run cmd/seed-users/main.go
```

---

## API Documentation

Swagger UI is available at:

```
http://localhost:8000/swagger/index.html
```

To regenerate docs after changing annotations:

```bash
swag init -g cmd/server/main.go
```
