package reporting

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/domain/attendance"
	"github.com/protone/erp/internal/domain/payroll"
)

type ReportingUseCase struct {
	payrollRepo    payroll.Repository
	attendanceRepo attendance.QueryRepository
}

func NewReportingUseCase(pr payroll.Repository, ar attendance.QueryRepository) *ReportingUseCase {
	return &ReportingUseCase{
		payrollRepo:    pr,
		attendanceRepo: ar,
	}
}

type AttendanceStats struct {
	TotalEmployees int `json:"total_employees"`
	Present        int `json:"present"`
	Late           int `json:"late"`
	OnLeave        int `json:"on_leave"`
	OvertimeToday  int `json:"overtime_today"`
}

func (uc *ReportingUseCase) GetDashboardData(ctx context.Context, tenantID uuid.UUID, month, year int) (map[string]interface{}, error) {
	// 1. Salary Report
	period, _ := uc.payrollRepo.GetPeriod(ctx, tenantID, month, year)
	salaryTotal := 0.0
	status := "N/A"
	if period != nil {
		salaryTotal = period.TotalAmount
		status = string(period.Status)
	}

	// 2. Department Breakdown
	breakdown, _ := uc.payrollRepo.GetDepartmentBreakdown(ctx, tenantID, month, year)

	// 3. Attendance Stats (Today)
	now := time.Now()
	todayStr := now.Format("2006-01-02")
	
	total, present, late, onLeave, _ := uc.attendanceRepo.GetStats(ctx, tenantID, todayStr)

	stats := AttendanceStats{
		TotalEmployees: total,
		Present:        present,
		Late:           late,
		OnLeave:        onLeave,
	}

	return map[string]interface{}{
		"salary": map[string]interface{}{
			"total_amount":  salaryTotal,
			"period_status": status,
			"breakdown":     breakdown,
		},
		"attendance_today": stats,
		"date_today":       todayStr,
	}, nil
}
