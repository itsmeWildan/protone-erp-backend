package employee

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ─── Value Objects ──────────────────────────────────────────────────────────

type Status string

const (
	StatusActive     Status = "active"
	StatusInactive   Status = "inactive"
	StatusTerminated Status = "terminated"
	StatusOnLeave    Status = "on_leave"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type EmploymentType string

const (
	EmploymentPermanent  EmploymentType = "permanent"
	EmploymentContract   EmploymentType = "contract"
	EmploymentInternship EmploymentType = "internship"
	EmploymentFreelance  EmploymentType = "freelance"
)

// ─── Entity ─────────────────────────────────────────────────────────────────

type Employee struct {
	ID              uuid.UUID       `json:"ID"`
	TenantID        uuid.UUID       `json:"TenantID"`
	UserID          *uuid.UUID      `json:"UserID"`
	NIK             string          `json:"NIK"`
	FullName        string          `json:"FullName"`
	Email           string          `json:"Email"`
	Phone           string          `json:"Phone"`
	Gender          Gender          `json:"Gender"`
	BirthDate       *time.Time      `json:"BirthDate"`
	DepartmentID    uuid.UUID       `json:"DepartmentID"`
	PositionID      uuid.UUID       `json:"PositionID"`
	ManagerID       *uuid.UUID      `json:"ManagerID"`
	EmploymentType  EmploymentType  `json:"EmploymentType"`
	Status          Status          `json:"Status"`
	JoinDate        time.Time       `json:"JoinDate"`
	EndDate         *time.Time      `json:"EndDate"`
	BasicSalary     float64         `json:"BasicSalary"`
	BankName        string          `json:"BankName"`
	BankAccountNo   string          `json:"BankAccountNo"`
	BankAccountName string          `json:"BankAccountName"`
	CreatedAt       time.Time       `json:"CreatedAt"`
	UpdatedAt       time.Time       `json:"UpdatedAt"`
}

// ─── Factory ────────────────────────────────────────────────────────────────

type CreateParams struct {
	TenantID       uuid.UUID
	NIK            string
	FullName       string
	Email          string
	Phone          string
	Gender         Gender
	BirthDate      *time.Time
	DepartmentID   uuid.UUID
	PositionID     uuid.UUID
	ManagerID      *uuid.UUID
	EmploymentType EmploymentType
	JoinDate       time.Time
	BasicSalary    float64
}

func NewEmployee(p CreateParams) (*Employee, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	now := time.Now()
	return &Employee{
		ID:             uuid.New(),
		TenantID:       p.TenantID,
		NIK:            p.NIK,
		FullName:       p.FullName,
		Email:          p.Email,
		Phone:          p.Phone,
		Gender:         p.Gender,
		BirthDate:      p.BirthDate,
		DepartmentID:   p.DepartmentID,
		PositionID:     p.PositionID,
		ManagerID:      p.ManagerID,
		EmploymentType: p.EmploymentType,
		Status:         StatusActive,
		JoinDate:       p.JoinDate,
		BasicSalary:    p.BasicSalary,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (p CreateParams) validate() error {
	if p.NIK == "" {
		return ErrNIKRequired
	}
	if p.FullName == "" {
		return ErrFullNameRequired
	}
	if p.TenantID == uuid.Nil {
		return ErrTenantIDRequired
	}
	if p.DepartmentID == uuid.Nil {
		return ErrDepartmentRequired
	}
	if p.PositionID == uuid.Nil {
		return ErrPositionRequired
	}
	if p.JoinDate.IsZero() {
		return ErrJoinDateRequired
	}
	return nil
}

// ─── Domain Behaviors ───────────────────────────────────────────────────────

func (e *Employee) Deactivate() error {
	if e.Status == StatusTerminated {
		return ErrAlreadyTerminated
	}
	e.Status = StatusInactive
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Employee) Terminate(endDate time.Time) error {
	if e.Status == StatusTerminated {
		return ErrAlreadyTerminated
	}
	e.Status = StatusTerminated
	e.EndDate = &endDate
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Employee) ChangePosition(departmentID, positionID uuid.UUID) error {
	if departmentID == uuid.Nil || positionID == uuid.Nil {
		return errors.New("department and position are required")
	}
	e.DepartmentID = departmentID
	e.PositionID = positionID
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Employee) UpdateBasicSalary(amount float64) error {
	if amount < 0 {
		return ErrInvalidSalary
	}
	e.BasicSalary = amount
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Employee) IsActive() bool {
	return e.Status == StatusActive
}
