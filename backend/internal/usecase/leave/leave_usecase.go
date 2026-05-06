package leave

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/employee"
	"github.com/protone/erp/internal/domain/leave"
)

type RequestLeaveInput struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	LeaveTypeID uuid.UUID
	StartDate   time.Time
	EndDate     time.Time
	Reason      string
}

type LeaveUseCase struct {
	leaveRepo    leave.Repository
	employeeRepo employee.Repository
}

func NewLeaveUseCase(lr leave.Repository, er employee.Repository) *LeaveUseCase {
	return &LeaveUseCase{
		leaveRepo:    lr,
		employeeRepo: er,
	}
}

func (uc *LeaveUseCase) GetLeaveTypes(ctx context.Context, tenantID uuid.UUID) ([]leave.LeaveType, error) {
	return uc.leaveRepo.GetTypes(ctx, tenantID)
}

func (uc *LeaveUseCase) RequestLeave(ctx context.Context, input RequestLeaveInput) (uuid.UUID, error) {
	// 1. Get employee by UserID
	emp, err := uc.employeeRepo.FindByUserID(ctx, input.TenantID, input.UserID)
	if err != nil {
		return uuid.Nil, err
	}
	if emp == nil {
		return uuid.Nil, errors.New("employee not found")
	}

	// 2. Calculate days (simple diff)
	days := int(input.EndDate.Sub(input.StartDate).Hours()/24) + 1
	if days <= 0 {
		return uuid.Nil, errors.New("invalid date range")
	}

	// 3. Check Balance
	year := input.StartDate.Year()
	balance, err := uc.leaveRepo.GetBalance(ctx, emp.ID, input.LeaveTypeID, year)
	if err != nil {
		return uuid.Nil, err
	}
	if balance == nil {
		// Auto-init balance for testing convenience
		balance = &leave.LeaveBalance{
			ID:          uuid.New(),
			TenantID:    input.TenantID,
			EmployeeID:  emp.ID,
			LeaveTypeID: input.LeaveTypeID,
			Year:        year,
			TotalDays:   12,
			UsedDays:    0,
			PendingDays: 0,
		}
		if err := uc.leaveRepo.CreateBalance(ctx, balance); err != nil {
			return uuid.Nil, err
		}
	}

	available := balance.TotalDays - balance.UsedDays - balance.PendingDays
	if available < days {
		return uuid.Nil, leave.ErrInsufficientBalance
	}

	// 4. Create Request
	req := leave.NewLeaveRequest(input.TenantID, emp.ID, input.LeaveTypeID, input.StartDate, input.EndDate, days, input.Reason)

	// 5. Update Balance (add to pending)
	balance.PendingDays += days
	
	// TODO: Use transaction if possible, but for now we do sequential
	if err := uc.leaveRepo.Save(ctx, req); err != nil {
		return uuid.Nil, err
	}

	return req.ID, uc.leaveRepo.UpdateBalance(ctx, balance)
}

func (uc *LeaveUseCase) ApproveLeave(ctx context.Context, tenantID, managerUserID, leaveID uuid.UUID) error {
	// 1. Get manager employee record
	mgr, err := uc.employeeRepo.FindByUserID(ctx, tenantID, managerUserID)
	if err != nil {
		return err
	}

	// 2. Get leave request
	l, err := uc.leaveRepo.GetByID(ctx, leaveID)
	if err != nil {
		return err
	}

	// 3. Approve
	if err := l.Approve(mgr.ID); err != nil {
		return err
	}

	// 4. Update Balance (move from pending to used)
	balance, err := uc.leaveRepo.GetBalance(ctx, l.EmployeeID, l.LeaveTypeID, l.StartDate.Year())
	if err != nil {
		return err
	}
	balance.PendingDays -= l.TotalDays
	balance.UsedDays += l.TotalDays

	if err := uc.leaveRepo.Update(ctx, l); err != nil {
		return err
	}

	return uc.leaveRepo.UpdateBalance(ctx, balance)
}

func (uc *LeaveUseCase) GetMyLeaves(ctx context.Context, tenantID, userID uuid.UUID) ([]leave.Leave, error) {
	emp, err := uc.employeeRepo.FindByUserID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("employee not found")
	}

	return uc.leaveRepo.GetByEmployee(ctx, emp.ID)
}
func (uc *LeaveUseCase) GiveBalance(ctx context.Context, tenantID, employeeID, leaveTypeID uuid.UUID, year, days int) error {
	balance := &leave.LeaveBalance{
		ID:          uuid.New(),
		TenantID:    tenantID,
		EmployeeID:  employeeID,
		LeaveTypeID: leaveTypeID,
		Year:        year,
		TotalDays:   days,
		UsedDays:    0,
		PendingDays: 0,
	}
	return uc.leaveRepo.CreateBalance(ctx, balance)
}
