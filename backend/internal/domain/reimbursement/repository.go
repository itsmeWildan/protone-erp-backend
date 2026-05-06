package reimbursement

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, r *Reimbursement) error
	GetByID(ctx context.Context, id uuid.UUID) (*Reimbursement, error)
	Update(ctx context.Context, r *Reimbursement) error
	GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]Reimbursement, error)
}
