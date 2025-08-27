# Go REST API Boilerplate

Clean Architecture boilerplate for building REST APIs using **Golang**.  
This project is designed to be scalable, maintainable, and easy to extend.

---

## Features

- [x] Authentication (Register & Login with JWT)
- [x] User Management (User CRUD)
- [x] Role Management
- [x] Middleware (Auth, Role Checker, Rate Limiter, CORS)
- [x] Database Migration & Seeding (PostgreSQL / MySQL)
- [x] Clean Architecture (separation of concerns)

---

## Tech Stack & Dependencies

**Core**

- Go 1.24.6
- Gin v1.10.1 – HTTP web framework
- GORM v1.30.1 – ORM for Go
- PostgreSQL & MySQL driver

**Security & Auth**

- JWT Authentication
- x/crypto – cryptography

**Validation & Internationalization**

- go-playground/validator/v10
- go-playground/locales
- go-playground/universal-translator

**Utilities**

- google/uuid – UUID generator
- joho/godotenv – Load `.env` file
- go-email-normalizer – Email normalization
- ulule/limiter – Rate limiter middleware

---

## Project Structure

```

go-rest-api-boilerplate/
├── cmd/                 # entry point (main.go)
├── internal/
│   ├── config/          # configuration (.env, env loader)
│   ├── domain/          # entities & repository interfaces
│   │   ├── entities/    # model entities (User, Role, AuthSession)
│   │   └── repositories # repository contracts
│   ├── infrastructure/  # technical implementations (db, repos)
│   │   ├── databases/   # migrations, seeding, DB connections
│   │   └── repositories # repository implementations
│   └── interfaces/      # delivery layer (HTTP)
│       └── http/
│           ├── controllers # endpoint handlers
│           ├── dto         # request/response DTOs
│           ├── mapper      # DTO <-> entity mapping
│           ├── middleware  # auth, cors, rate limiter, etc.
│           └── router      # HTTP routing
├── pkg/                  # helpers/utilities
├── tmp/                  # runtime temp (e.g., logs)
├── api/                  # Postman collection
│   └── go-rest-api-boilerplate.postman\_collection.json
├── .air.toml             # Air live reload config
├── .env.example          # environment example
├── go.mod
├── CONTRIBUTING.md
├── CHANGELOG.md
└── README.md

````

---

## API Endpoints

### API Documentation

This project comes with a **Postman Collection**.  
Available here: [Postman Collection](./api/go-rest-api-boilerplate.postman_collection.json)

### Auth

| Method | Endpoint         | Description   |
|--------|------------------|---------------|
| POST   | `/auth/register` | Register user |
| POST   | `/auth/login`    | Login user    |

### Users

| Method | Endpoint      | Description       |
|--------|---------------|-------------------|
| POST   | `/users`      | Create User       |
| PATCH  | `/users/{id}` | Update User by ID |
| DELETE | `/users/{id}` | Delete User by ID |
| GET    | `/users/{id}` | Get User by ID    |
| GET    | `/users`      | Get All Users     |

### Roles

| Method | Endpoint | Description   |
|--------|----------|---------------|
| GET    | `/roles` | Get All Roles |

---

## Installation

Clone repo:

```bash
git clone https://github.com/username/go-rest-api-boilerplate.git
cd go-rest-api-boilerplate
````

Install dependencies:

```bash
go mod tidy
```

Copy env:

```bash
cp .env.example .env
```

Run migrations & seed:

```bash
go run internal/infrastructure/databases/migrate.go
go run internal/infrastructure/databases/seed.go
```

---

## Running

With Air (auto reload):

```bash
air
```

Manual:

```bash
go run cmd/main.go
```

API will run at `http://localhost:3000/api/v1`

---

## Tests

```bash
go test ./...
```

---

## Contributing

See the contribution guide in [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## Changelog

All version changes are documented in [CHANGELOG.md](./CHANGELOG.md).

---

## License

MIT License © Alfin Noviaji