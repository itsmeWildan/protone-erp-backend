package overtime

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, o *OvertimeRequest) error
	Update(ctx context.Context, o *OvertimeRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*OvertimeRequest, error)
	GetByEmployeeAndMonth(ctx context.Context, employeeID uuid.UUID, month, year int) ([]OvertimeRequest, error)
	GetApprovedSumByEmployee(ctx context.Context, employeeID uuid.UUID, month, year int) (float64, error)
}
