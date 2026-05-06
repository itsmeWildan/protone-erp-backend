package employee

import (
	"context"

	"github.com/google/uuid"
)

// ─── Write Port (Command Repository) ────────────────────────────────────────
// Digunakan oleh UseCase command (tulis data pakai entity)

type Repository interface {
	Save(ctx context.Context, e *Employee) error
	Update(ctx context.Context, e *Employee) error
	Delete(ctx context.Context, tenantID, id uuid.UUID) error
	FindByID(ctx context.Context, tenantID, id uuid.UUID) (*Employee, error)
	FindByNIK(ctx context.Context, tenantID uuid.UUID, nik string) (*Employee, error)
	FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*Employee, error)
}

// ─── Read Port (Query Repository) ───────────────────────────────────────────
// Digunakan oleh UseCase query (baca data, bisa pakai JOIN / flat DTO)
// semi-CQRS: read pakai model yang berbeda dari entity

type Filter struct {
	TenantID     uuid.UUID
	DepartmentID *uuid.UUID
	PositionID   *uuid.UUID
	Status       *Status
	Search       string // cari di NIK / nama
	Page         int
	PerPage      int
}

// EmployeeListItem — read model untuk list (JOIN dengan dept & position name)
type EmployeeListItem struct {
	ID             string
	NIK            string
	FullName       string
	Email          string
	DepartmentName string
	PositionName   string
	Status         string
	JoinDate       string
	EmploymentType string
}

// EmployeeDetail — read model lengkap untuk detail view
type EmployeeDetail struct {
	ID              string
	NIK             string
	FullName        string
	Email           string
	Phone           string
	Gender          string
	BirthDate       *string
	DepartmentID    string
	DepartmentName  string
	PositionID      string
	PositionName    string
	ManagerID       *string
	ManagerName     *string
	EmploymentType  string
	Status          string
	JoinDate        string
	BasicSalary     float64
	BankName        string
	BankAccountNo   string
	BankAccountName string
}

type QueryRepository interface {
	List(ctx context.Context, filter Filter) ([]EmployeeListItem, int64, error)
	GetByID(ctx context.Context, tenantID, id uuid.UUID) (*EmployeeDetail, error)
}
