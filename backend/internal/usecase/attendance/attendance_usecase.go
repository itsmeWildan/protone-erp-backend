package attendance

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/attendance"
	"github.com/protone/erp/internal/domain/employee"
)

type UseCase interface {
	ClockIn(ctx context.Context, input ClockInInput) error
	ClockOut(ctx context.Context, input ClockOutInput) error
}

type ClockInInput struct {
	TenantID   uuid.UUID
	UserID     uuid.UUID
	LocationIn string
	Notes      string
}

type ClockOutInput struct {
	TenantID    uuid.UUID
	UserID      uuid.UUID
	LocationOut string
}

type AttendanceUseCase struct {
	attendanceRepo attendance.Repository
	employeeRepo   employee.Repository
}

func NewAttendanceUseCase(ar attendance.Repository, er employee.Repository) *AttendanceUseCase {
	return &AttendanceUseCase{
		attendanceRepo: ar,
		employeeRepo:   er,
	}
}

func (uc *AttendanceUseCase) ClockIn(ctx context.Context, input ClockInInput) error {
	// 1. Get employee by UserID
	emp, err := uc.employeeRepo.FindByUserID(ctx, input.TenantID, input.UserID)
	if err != nil {
		return err
	}
	if emp == nil {
		return errors.New("employee record not found for this user")
	}

	// 2. Check if already clocked in today
	today := time.Now()
	existing, err := uc.attendanceRepo.GetByEmployeeAndDate(ctx, emp.ID, today)
	if err != nil {
		return err
	}
	if existing != nil {
		return attendance.ErrAlreadyClockedIn
	}

	// 3. Create new attendance
	attr := attendance.NewAttendance(input.TenantID, emp.ID, input.LocationIn, input.Notes)
	
	return uc.attendanceRepo.Save(ctx, attr)
}

func (uc *AttendanceUseCase) ClockOut(ctx context.Context, input ClockOutInput) error {
	// 1. Get employee by UserID
	emp, err := uc.employeeRepo.FindByUserID(ctx, input.TenantID, input.UserID)
	if err != nil {
		return err
	}
	if emp == nil {
		return errors.New("employee record not found for this user")
	}

	// 2. Get today's attendance
	today := time.Now()
	attr, err := uc.attendanceRepo.GetByEmployeeAndDate(ctx, emp.ID, today)
	if err != nil {
		return err
	}
	if attr == nil {
		return attendance.ErrNotClockedIn
	}

	// 3. Clock out
	if err := attr.ClockOut(input.LocationOut); err != nil {
		return err
	}

	return uc.attendanceRepo.Update(ctx, attr)
}
