# PROTONE ERP Backend

PROTONE ERP is a Go (Golang) based Enterprise Resource Planning (ERP) system designed with Clean Architecture and Domain-Driven Design (DDD) principles. This backend handles the entire core business logic for Human Resources (HR) and Finance management.

## 🚀 Key Features

- **Multi-Tenancy**: Supports multiple companies (tenants) within a single database instance.
- **Employee Management**: Manage employee master data, positions, and departments.
- **Attendance**: Clock-In/Out features with location tracking and notes.
- **Leave Management**: Leave requests, manager approvals, and annual leave balance tracking.
- **Overtime**: Overtime submissions with automatic calculation into payroll.
- **Reimbursement**: Operational expense claims with an approval system.
- **Payroll**: 
  - Automated monthly payroll generation.
  - Calculation of basic salary, allowances, deductions, and overtime pay.
  - Export payroll slips to PDF format.
  - Automatic integration with Financial Journals (Accounting).
- **Finance & Budgeting**: 
  - Departmental Budget Management.
  - Automatic General Ledger (GL) entry recording from payroll transactions.

## 🛠️ Technology Stack

- **Language**: Go (Golang) 1.21+
- **Database**: PostgreSQL (pgx driver)
- **HTTP Framework**: Chi Router 
- **Library**: 
  - `google/uuid` (Identity)
  - `maroto` (PDF Generation)
  - `golang-jwt` (Authentication)
  - `migrate` (Database Migrations)

## 🏗️ Project Structure (Clean Architecture)

- `internal/domain`: Business entities and core logic (Zero external dependencies).
- `internal/usecase`: Application workflows (Orchestration logic).
- `internal/infrastructure`: Database implementation (Postgres), Migrations, and External Tools.
- `internal/delivery/http`: API Handlers, Router, and Middleware.
- `pkg/`: Helper libraries for reusability (Response wrapper, JWT, PDF).

## 🏃 Getting Started

### 1. Prerequisites
- Instal [Go](https://go.dev/dl/)
- Instal [PostgreSQL](https://www.postgresql.org/download/)

### 2. Database Configuration
Create a database named `protone_erp` in your PostgreSQL instance.

### 3. Environment Setup
Create a `.env` file in the `backend/` folder (example content below):
```env
DB_URL=postgres://user:password@localhost:5432/protone_erp?sslmode=disable
JWT_SECRET=your_secret_key
PORT=3000
```

### 4. Running Migrations & Seeders
```bash
# Run migrations to create tables
# (Or use a migration tool if available)

# Run seeder for initial data (Tenant BCDE)
go run seed_payroll.go
```

### 5. Running the Application
```bash
cd backend
go run cmd/api/main.go
```

## 🧪 API Testing
You can test the APIs using Postman or HTTPie. API collections are available in the docs/ folder (optional).

---
Maintained by itsmeWildan © 2026 PROTONE Project.


