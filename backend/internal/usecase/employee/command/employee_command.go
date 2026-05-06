package command

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
)

// ─── Create Employee ─────────────────────────────────────────────────────────

type CreateEmployeeInput struct {
	TenantID       uuid.UUID
	NIK            string
	FullName       string
	Email          string
	Phone          string
	Gender         string
	BirthDate      *time.Time
	DepartmentID   uuid.UUID
	PositionID     uuid.UUID
	ManagerID      *uuid.UUID
	EmploymentType string
	JoinDate       time.Time
	BasicSalary    float64
}

type CreateEmployeeOutput struct {
	ID string
}

type CreateEmployeeUseCase struct {
	repo employee.Repository
}

func NewCreateEmployee(repo employee.Repository) *CreateEmployeeUseCase {
	return &CreateEmployeeUseCase{repo: repo}
}

func (uc *CreateEmployeeUseCase) Execute(ctx context.Context, input CreateEmployeeInput) (*CreateEmployeeOutput, error) {
	// Check NIK uniqueness
	existing, err := uc.repo.FindByNIK(ctx, input.TenantID, input.NIK)
	if err != nil {
		return nil, fmt.Errorf("check NIK: %w", err)
	}
	if existing != nil {
		return nil, employee.ErrNIKAlreadyExists
	}

	// Build entity via factory
	e, err := employee.NewEmployee(employee.CreateParams{
		TenantID:       input.TenantID,
		NIK:            input.NIK,
		FullName:       input.FullName,
		Email:          input.Email,
		Phone:          input.Phone,
		Gender:         employee.Gender(input.Gender),
		BirthDate:      input.BirthDate,
		DepartmentID:   input.DepartmentID,
		PositionID:     input.PositionID,
		ManagerID:      input.ManagerID,
		EmploymentType: employee.EmploymentType(input.EmploymentType),
		JoinDate:       input.JoinDate,
		BasicSalary:    input.BasicSalary,
	})
	if err != nil {
		return nil, fmt.Errorf("create employee entity: %w", err)
	}

	// Persist
	if err := uc.repo.Save(ctx, e); err != nil {
		return nil, fmt.Errorf("save employee: %w", err)
	}

	return &CreateEmployeeOutput{ID: e.ID.String()}, nil
}

// ─── Update Employee ─────────────────────────────────────────────────────────

type UpdateEmployeeInput struct {
	TenantID    uuid.UUID
	ID          uuid.UUID
	FullName    string
	Email       string
	Phone       string
	DepartmentID uuid.UUID
	PositionID  uuid.UUID
	ManagerID   *uuid.UUID
	BasicSalary float64
	BankName    string
	BankAccountNo   string
	BankAccountName string
}

type UpdateEmployeeUseCase struct {
	repo employee.Repository
}

func NewUpdateEmployee(repo employee.Repository) *UpdateEmployeeUseCase {
	return &UpdateEmployeeUseCase{repo: repo}
}

func (uc *UpdateEmployeeUseCase) Execute(ctx context.Context, input UpdateEmployeeInput) error {
	e, err := uc.repo.FindByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return err
	}

	// Apply changes via domain behaviors
	if err := e.ChangePosition(input.DepartmentID, input.PositionID); err != nil {
		return err
	}
	if err := e.UpdateBasicSalary(input.BasicSalary); err != nil {
		return err
	}

	// Update plain fields
	e.FullName = input.FullName
	e.Email = input.Email
	e.Phone = input.Phone
	e.ManagerID = input.ManagerID
	e.BankName = input.BankName
	e.BankAccountNo = input.BankAccountNo
	e.BankAccountName = input.BankAccountName

	return uc.repo.Update(ctx, e)
}

// ─── Delete Employee ─────────────────────────────────────────────────────────

type DeleteEmployeeUseCase struct {
	repo employee.Repository
}

func NewDeleteEmployee(repo employee.Repository) *DeleteEmployeeUseCase {
	return &DeleteEmployeeUseCase{repo: repo}
}

func (uc *DeleteEmployeeUseCase) Execute(ctx context.Context, tenantID, id uuid.UUID) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
