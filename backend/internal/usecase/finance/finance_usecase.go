package finance

import (
	"context"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/finance"
)

type UseCase interface {
	SetBudget(ctx context.Context, budget *finance.DepartmentBudget) error
	CheckBudget(ctx context.Context, tenantID, deptID uuid.UUID, month, year int) (*finance.DepartmentBudget, error)
}

type financeUseCase struct {
	repo finance.Repository
}

func NewUseCase(repo finance.Repository) UseCase {
	return &financeUseCase{repo: repo}
}

func (u *financeUseCase) SetBudget(ctx context.Context, budget *finance.DepartmentBudget) error {
	return u.repo.UpdateBudget(ctx, budget)
}

func (u *financeUseCase) CheckBudget(ctx context.Context, tenantID, deptID uuid.UUID, month, year int) (*finance.DepartmentBudget, error) {
	return u.repo.GetBudget(ctx, tenantID, deptID, month, year)
}
