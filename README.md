# Golang Fiber Boilerplate

Production ready backend boilerplate built with Go Fiber, Ent ORM, PostgreSQL, Redis, and Atlas Migration. Includes complete authentication system, user profile, and common middleware to accelerate your project development.

![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?logo=go)
![Fiber](https://img.shields.io/badge/Fiber-v2.x-00ADD8?logo=fiber)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15%2B-336791?logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Available Features

- ✅ JWT Authentication (Register, Login, Logout)
- ✅ User Profile Management
- ✅ Authentication Middleware
- ✅ Standardized Error Handler
- ✅ Rate Limiter
- ✅ Request Logger
- ✅ Ent ORM Database Layer
- ✅ Atlas Database Migration
- ✅ Hot Reload with Air
- ✅ Redis for Session & Cache
- ✅ AWS S3 Storage Integration
- ✅ Mailer Service
- ✅ Well Organized Folder Structure

---

## 📋 System Requirements

Make sure you have these installed on your system:

- Go 1.24 or newer
- PostgreSQL 15+
- Redis 7+
- Docker (optional but recommended)

---

## 🚀 Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/tfajar123/go-boilerplate.git
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

4.  Edit `.env` file and adjust database, redis and other configurations accordingly.

5.  Install required development tools:

    ```bash
    # Install Atlas for database migration
    curl -sSf https://atlasgo.sh | sh

    # Install Air for hot reload
    go install github.com/air-verse/air@latest
    ```

---

## ⚡ Running The Application

For development with hot reload:

```bash
air
```

Application will run on `http://localhost:3000` by default.

---

## 🛠️ Makefile Usage

This project includes `Makefile` to simplify running common commands. All commands are executed using `make <command_name>` format.

| Command                                 | Description                                                               |
| --------------------------------------- | ------------------------------------------------------------------------- |
| `make help`                             | Display all available commands                                            |
| `make setup`                            | First time project setup (generate ent, hash migration, apply migrations) |
| `make gen`                              | Regenerate Ent ORM client after schema changes                            |
| `make migrate-hash`                     | Generate integrity hash for migration files                               |
| `make migrate-diff name=migration_name` | Create new migration file based on schema changes                         |
| `make migrate-apply`                    | Run all pending migrations                                                |
| `make migrate-local`                    | Apply migrations to local environment                                     |
| `make migrate-staging`                  | Apply migrations to staging environment                                   |
| `make migrate-prod`                     | Apply migrations to production environment                                |

### Makefile Workflow Example:

```bash
# After modifying schema in ent/schema/
make gen

# Create migration file
make migrate-diff name=add_address_column

# Execute migration
make migrate-apply
```

> 💡 Tip: Use `make setup` when first setting up the project, this will automatically run all required steps for you.

---

## 📖 Migration Guide

This project uses **Schema First** approach:

1.  Edit or add schema files in `./ent/schema/` directory
2.  Generate Ent client: `make gen`
3.  Create migration: `make migrate-diff name=change_description`
4.  Review generated migration file in `./migrations/` folder
5.  Run migration: `make migrate-apply`

For other environments use:

```bash
make migrate-staging
make migrate-prod
```

---

## 📂 Folder Structure

```
go-boilerplate/
├── apps/
│   ├── cmd/server/          # Application entry point
│   └── internal/
│       ├── config/          # Environment config loader
│       ├── database/        # Database & Redis connection
│       ├── features/        # Application features (modular)
│       ├── middleware/      # Global middleware
│       ├── route/           # Routing definitions
│       └── utils/           # Helpers & common functions
├── ent/                     # Ent ORM schema & generated code
├── migrations/              # SQL Migration files
├── .env.example             # Example environment configuration
├── makefile                 # Make commands
├── atlas.hcl                # Atlas migration configuration
└── docker-compose.yml       # Development docker stack
```

---

## 🐳 Using Docker

To run all services (PostgreSQL, Redis, App) with single command:

```bash
docker compose up -d
```

---

## 🤝 Contributing

For major changes, please open an issue first to discuss what you would like to change. Always make sure to update tests as appropriate.

---

## 📄 License

This project is licensed under **MIT** License. See [LICENSE](LICENSE) file for full details.
