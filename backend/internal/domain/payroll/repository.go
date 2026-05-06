package payroll

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetComponents(ctx context.Context, tenantID uuid.UUID) ([]SalaryComponent, error)
	GetPeriodByID(ctx context.Context, id uuid.UUID) (*PayrollPeriod, error)
	GetPeriod(ctx context.Context, tenantID uuid.UUID, month, year int) (*PayrollPeriod, error)
	CreatePeriod(ctx context.Context, p *PayrollPeriod) error
	UpdatePeriod(ctx context.Context, p *PayrollPeriod) error
	UpdatePeriodStatus(ctx context.Context, id uuid.UUID, status string) error
	GetSlipsByPeriod(ctx context.Context, periodID uuid.UUID) ([]PayrollSlip, error)
	GetSlipByEmployee(ctx context.Context, periodID, employeeID uuid.UUID) (*PayrollSlip, error)
	SaveSlip(ctx context.Context, s *PayrollSlip) error
	DeleteSlipsByPeriod(ctx context.Context, periodID uuid.UUID) error
	GetDepartmentBreakdown(ctx context.Context, tenantID uuid.UUID, month, year int) (map[uuid.UUID]float64, error)
}
