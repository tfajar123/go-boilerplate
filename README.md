# Golang Fiber Boilerplate (MongoDB Variant)

Production ready backend boilerplate built with Go Fiber, MongoDB, Redis, and MinIO S3. Includes complete authentication system, user profile, and common middleware to accelerate your project development.

![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2.x-00ADD8?logo=fiber)
![MongoDB](https://img.shields.io/badge/MongoDB-7%2B-47A248?logo=mongodb)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Available Features

- ✅ JWT Authentication (Register, Login, Logout)
- ✅ OTP Email Verification
- ✅ Forgot & Reset Password
- ✅ User Profile Management
- ✅ Authentication Middleware
- ✅ Standardized Error Handler
- ✅ Rate Limiter
- ✅ Request Logger
- ✅ Schema-First MongoDB Models
- ✅ Code-Based Index Management
- ✅ Hot Reload with Air
- ✅ Redis for Session & Cache
- ✅ AWS S3 / MinIO Storage Integration
- ✅ Mailer Service
- ✅ Well Organized Folder Structure
- ✅ Docker & Docker Compose Ready
- ✅ Jenkins CI/CD Pipeline

---

## 📋 System Requirements

Make sure you have these installed on your system:

- Go 1.25 or newer
- MongoDB 7+
- Redis 7+
- Docker (optional but recommended)

---

## 🚀 Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/tfajar123/go-boilerplate.git -b var/mongodb
    cd go-boilerplate
    ```

2.  Install Go dependencies:

    ```bash
    go mod download
    go mod tidy
    ```

3.  Copy environment configuration file:

    ```bash
    cp .env.example .env
    ```

4.  Edit `.env` file and adjust MongoDB, Redis and other configurations:

    ```env
    MONGO_URI=mongodb://mongoadmin:mongo123@localhost:27017
    MONGO_DB_NAME=go_boilerplate
    ```

5.  Install Air for hot reload (optional):

    ```bash
    go install github.com/air-verse/air@latest
    ```

---

## ⚡ Running The Application

Start MongoDB and Redis first (if not already running):

```bash
docker compose --profile dev up -d
```

Then run the application with hot reload:

```bash
air
```

Or without Air:

```bash
go run apps/cmd/server/main.go
```

Application will run on `http://localhost:3098` by default.

---

## 🛠️ Makefile Usage

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `make help`  | Display all available commands                             |
| `make setup` | First time project setup (install deps)                    |
| `make dev`   | Start development server with Air                          |

---

## 📖 Schema & Migration Guide

This project uses a **Schema-First** approach with Go structs as the single source of truth for MongoDB collections.

### How It Works

| Component | File | Purpose |
|---|---|---|
| Schema | `models/user.go` | Defines collection structure with BSON tags |
| Indexes | `models/indexes.go` | Declares all MongoDB indexes |

### Adding a New Collection

1.  Create a new model file in `models/`:

    ```go
    // models/product.go
    type Product struct {
        ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
        Name        string        `bson:"name" json:"name"`
        Price       float64       `bson:"price" json:"price"`
        CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
        UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
    }

    const CollectionProducts = "products"
    ```

2.  Add indexes in `models/indexes.go` (if needed):

    ```go
    productsIndexes := []mongo.IndexModel{
        {
            Keys: bson.D{{Key: "name", Value: "text"}},
        },
    }
    _, err = db.Collection(CollectionProducts).Indexes().CreateMany(ctx, productsIndexes)
    ```

3.  **Done!** No migration files needed. MongoDB auto-creates the collection on first insert, and indexes are applied on every server startup.

### Common Schema Changes

| Change | Action Required |
|---|---|
| Add a new field | Add field to struct → restart |
| Remove a field | Remove from struct → restart |
| Add a new collection | Create model file → use in service |
| Add an index | Add to `indexes.go` → restart |
| Rename a field | Requires a one-time data migration script |

> 💡 **Why no migration files?** MongoDB is schemaless — the database accepts any document structure. Your Go structs enforce data consistency at the application level, and indexes are managed declaratively via code.

---

## 📂 Folder Structure

```
go-boilerplate/
├── apps/
│   ├── cmd/server/          # Application entry point
│   └── internal/
│       ├── config/          # Environment config loader
│       ├── database/        # MongoDB, Redis & S3 connection
│       ├── features/        # Application features (modular)
│       │   ├── auth/        # Authentication (handlers, services, dto, validation)
│       │   ├── profile/     # User profile management
│       │   ├── mailer/      # Email service
│       │   └── storage/     # S3/MinIO file upload service
│       ├── middleware/      # Global middleware
│       ├── route/           # Routing definitions
│       └── utils/           # Helpers & common functions
├── models/                  # Schema-first MongoDB models & indexes
├── .env.example             # Example environment configuration
├── makefile                 # Make commands
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Development docker stack
└── docker-compose.staging.yml # Staging deployment
```

---

## 🐳 Using Docker

To run all services (MongoDB, Redis, App) with a single command:

```bash
# Development (includes MongoDB & Redis containers)
docker compose --profile dev up -d

# Production
docker compose up -d
```

---

## 🤝 Contributing

For major changes, please open an issue first to discuss what you would like to change. Always make sure to update tests as appropriate.

---

## 📄 License

This project is licensed under **MIT** License. See [LICENSE](LICENSE) file for full details.
