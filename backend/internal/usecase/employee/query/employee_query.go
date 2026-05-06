package query

import (
	"context"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
)

// ─── List Employees ──────────────────────────────────────────────────────────

type ListEmployeesInput struct {
	TenantID     uuid.UUID
	DepartmentID *uuid.UUID
	Status       *string
	Search       string
	Page         int
	PerPage      int
}

type ListEmployeesOutput struct {
	Items   []employee.EmployeeListItem
	Total   int64
	Page    int
	PerPage int
}

type ListEmployeesUseCase struct {
	queryRepo employee.QueryRepository
}

func NewListEmployees(queryRepo employee.QueryRepository) *ListEmployeesUseCase {
	return &ListEmployeesUseCase{queryRepo: queryRepo}
}

func (uc *ListEmployeesUseCase) Execute(ctx context.Context, input ListEmployeesInput) (*ListEmployeesOutput, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.PerPage <= 0 || input.PerPage > 100 {
		input.PerPage = 20
	}

	filter := employee.Filter{
		TenantID: input.TenantID,
		Search:   input.Search,
		Page:     input.Page,
		PerPage:  input.PerPage,
	}
	if input.DepartmentID != nil {
		filter.DepartmentID = input.DepartmentID
	}
	if input.Status != nil {
		s := employee.Status(*input.Status)
		filter.Status = &s
	}

	items, total, err := uc.queryRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &ListEmployeesOutput{
		Items:   items,
		Total:   total,
		Page:    input.Page,
		PerPage: input.PerPage,
	}, nil
}

// ─── Get Employee Detail ─────────────────────────────────────────────────────

type GetEmployeeUseCase struct {
	queryRepo employee.QueryRepository
}

func NewGetEmployee(queryRepo employee.QueryRepository) *GetEmployeeUseCase {
	return &GetEmployeeUseCase{queryRepo: queryRepo}
}

func (uc *GetEmployeeUseCase) Execute(ctx context.Context, tenantID, id uuid.UUID) (*employee.EmployeeDetail, error) {
	return uc.queryRepo.GetByID(ctx, tenantID, id)
}
