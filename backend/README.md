# ProtoERP Backend

REST API untuk Employee & Finance Management System.

**Stack**: Go · chi · pgx · PostgreSQL · JWT

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- `golang-migrate` CLI (untuk jalankan migrations)

## Quick Start

### 1. Install Go
Download dari https://go.dev/dl/ (minimal v1.22)

### 2. Jalankan PostgreSQL via Docker
```bash
docker compose up -d
```

### 3. Copy environment config
```bash
cp .env.example .env
# Edit .env sesuai kebutuhan
```

### 4. Install golang-migrate
```bash
# Windows (via scoop)
scoop install migrate

# Atau download langsung dari:
# https://github.com/golang-migrate/migrate/releases
```

### 5. Jalankan migrasi database
```bash
migrate -path internal/infrastructure/persistence/migrations \
        -database "postgres://postgres:postgres@localhost:5432/protone_erp?sslmode=disable" \
        up
```

### 6. Download Go dependencies
```bash
go mod tidy
```

### 7. Jalankan API server
```bash
go run ./cmd/api
```

API akan berjalan di http://localhost:8080

## Struktur Folder

```
backend/
├── cmd/api/             # Entry point
├── config/              # Configuration loader
├── internal/
│   ├── domain/          # Business entities, ports (interfaces)
│   │   └── employee/    # entity.go, repository.go, errors.go
│   ├── usecase/         # Application use cases
│   │   └── employee/
│   │       ├── command/ # CreateEmployee, UpdateEmployee, DeleteEmployee
│   │       └── query/   # ListEmployees, GetEmployee
│   ├── delivery/
│   │   └── http/        # HTTP handlers, middleware, router
│   └── infrastructure/
│       └── persistence/
│           ├── postgres/ # DB adapter implementations
│           └── migrations/ # SQL migration files
└── pkg/
    ├── jwt/             # JWT token manager
    └── response/        # Standard JSON response helpers
```

## API Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /health | Health check |
| GET | /api/v1/employees | List employees (paginated) |
| POST | /api/v1/employees | Create employee |
| GET | /api/v1/employees/:id | Get employee detail |
| PUT | /api/v1/employees/:id | Update employee |
| DELETE | /api/v1/employees/:id | Delete employee |

## Response Format

```json
{
  "status": "success",
  "message": "Employees fetched successfully",
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 150
  }
}
```

## Dev Tools

- pgAdmin: http://localhost:5050 (email: admin@protone.local / password: admin)
