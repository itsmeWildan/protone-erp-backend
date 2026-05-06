package overtime

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
	"github.com/protone/erp/internal/domain/overtime"
)

type SubmitOvertimeInput struct {
	TenantID      uuid.UUID
	UserID        uuid.UUID
	Date          time.Time
	StartTime     string
	EndTime       string
	DurationHours float64
	Reason        string
}

type OvertimeUseCase struct {
	repo         overtime.Repository
	employeeRepo employee.Repository
}

func NewOvertimeUseCase(r overtime.Repository, er employee.Repository) *OvertimeUseCase {
	return &OvertimeUseCase{
		repo:         r,
		employeeRepo: er,
	}
}

func (uc *OvertimeUseCase) SubmitRequest(ctx context.Context, input SubmitOvertimeInput) error {
	emp, err := uc.employeeRepo.FindByUserID(ctx, input.TenantID, input.UserID)
	if err != nil {
		return err
	}
	if emp == nil {
		return fmt.Errorf("employee not found")
	}

	req := overtime.NewOvertimeRequest(input.TenantID, emp.ID, input.Date, input.StartTime, input.EndTime, input.DurationHours, input.Reason)
	return uc.repo.Save(ctx, req)
}

func (uc *OvertimeUseCase) ApproveRequest(ctx context.Context, tenantID, managerUserID, requestID uuid.UUID) error {
	req, err := uc.repo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return fmt.Errorf("request not found")
	}

	req.Approve(managerUserID)
	return uc.repo.Update(ctx, req)
}

func (uc *OvertimeUseCase) GetMyOvertime(ctx context.Context, tenantID, userID uuid.UUID, month, year int) ([]overtime.OvertimeRequest, error) {
	emp, err := uc.employeeRepo.FindByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, fmt.Errorf("employee not found")
	}

	return uc.repo.GetByEmployeeAndMonth(ctx, emp.ID, month, year)
}
