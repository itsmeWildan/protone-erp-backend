package leave

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetTypes(ctx context.Context, tenantID uuid.UUID) ([]LeaveType, error)
	Save(ctx context.Context, l *Leave) error
	GetByID(ctx context.Context, id uuid.UUID) (*Leave, error)
	Update(ctx context.Context, l *Leave) error
	GetBalance(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*LeaveBalance, error)
	CreateBalance(ctx context.Context, b *LeaveBalance) error
	UpdateBalance(ctx context.Context, b *LeaveBalance) error
	GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]Leave, error)
}
