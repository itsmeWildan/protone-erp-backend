package payroll

import (
	"time"

	"github.com/google/uuid"
)

type ComponentType string

const (
	TypeAllowance ComponentType = "allowance"
	TypeDeduction ComponentType = "deduction"
)

const (
	StatusDraft      = "draft"
	StatusProcessing = "processing"
	StatusApproved   = "approved"
	StatusPaid       = "paid"
	StatusCancelled  = "cancelled"
)

type SalaryComponent struct {
	ID            uuid.UUID     `json:"ID"`
	TenantID      uuid.UUID     `json:"TenantID"`
	Name          string        `json:"Name"`
	Type          ComponentType `json:"Type"`
	Code          string        `json:"Code"`
	DefaultAmount float64       `json:"DefaultAmount"`
	IsTaxable     bool          `json:"IsTaxable"`
	CreatedAt     time.Time     `json:"CreatedAt"`
	UpdatedAt     time.Time     `json:"UpdatedAt"`
}

type PayrollPeriod struct {
	ID          uuid.UUID `json:"ID"`
	TenantID    uuid.UUID `json:"TenantID"`
	Month       int       `json:"Month"`
	Year        int       `json:"Year"`
	Status      string    `json:"Status"`
	TotalAmount float64   `json:"TotalAmount"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}

func NewPayrollPeriod(tenantID uuid.UUID, month, year int) *PayrollPeriod {
	return &PayrollPeriod{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Month:     month,
		Year:      year,
		Status:    StatusDraft,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type PayrollSlip struct {
	ID              uuid.UUID           `json:"ID"`
	TenantID        uuid.UUID           `json:"TenantID"`
	PayrollPeriodID uuid.UUID           `json:"PayrollPeriodID"`
	EmployeeID      uuid.UUID           `json:"EmployeeID"`
	BasicSalary     float64             `json:"BasicSalary"`
	TotalAllowance  float64             `json:"TotalAllowance"`
	TotalDeduction  float64             `json:"TotalDeduction"`
	NetSalary       float64             `json:"NetSalary"`
	WorkingDays     int                 `json:"WorkingDays"`
	PresentDays     int                 `json:"PresentDays"`
	OvertimeHours   float64             `json:"OvertimeHours"`
	OvertimeAmount  float64             `json:"OvertimeAmount"`
	Details         []PayrollSlipDetail `json:"Details"`
	CreatedAt       time.Time           `json:"CreatedAt"`
	UpdatedAt       time.Time           `json:"UpdatedAt"`
}

func (s *PayrollSlip) CalculateNet() {
	s.NetSalary = s.BasicSalary + s.TotalAllowance + s.OvertimeAmount - s.TotalDeduction
}

type PayrollSlipDetail struct {
	ID                uuid.UUID     `json:"ID"`
	PayrollSlipID     uuid.UUID     `json:"PayrollSlipID"`
	SalaryComponentID uuid.UUID     `json:"SalaryComponentID"`
	Type              ComponentType `json:"Type"`
	Amount            float64       `json:"Amount"`
}
