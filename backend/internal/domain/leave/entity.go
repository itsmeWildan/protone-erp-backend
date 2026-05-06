package leave

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInsufficientBalance = errors.New("insufficient leave balance")

type LeaveType struct {
	ID          uuid.UUID `json:"ID"`
	TenantID    uuid.UUID `json:"TenantID"`
	Name        string    `json:"Name"`
	Description string    `json:"Description"`
	MaxDays     int       `json:"MaxDays"`
	CreatedAt   time.Time `json:"CreatedAt"`
}

type Leave struct {
	ID           uuid.UUID  `json:"ID"`
	TenantID     uuid.UUID  `json:"TenantID"`
	EmployeeID   uuid.UUID  `json:"EmployeeID"`
	LeaveTypeID  uuid.UUID  `json:"LeaveTypeID"`
	StartDate    time.Time  `json:"StartDate"`
	EndDate      time.Time  `json:"EndDate"`
	TotalDays    int        `json:"TotalDays"`
	Reason       string     `json:"Reason"`
	Status       string     `json:"Status"`
	ApprovedBy   *uuid.UUID `json:"ApprovedBy"`
	ApprovedDate *time.Time `json:"ApprovedDate"`
	CreatedAt    time.Time  `json:"CreatedAt"`
	UpdatedAt    time.Time  `json:"UpdatedAt"`
}

func NewLeaveRequest(tenantID, empID, typeID uuid.UUID, start, end time.Time, days int, reason string) *Leave {
	return &Leave{
		ID:          uuid.New(),
		TenantID:    tenantID,
		EmployeeID:  empID,
		LeaveTypeID: typeID,
		StartDate:   start,
		EndDate:     end,
		TotalDays:   days,
		Reason:      reason,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (l *Leave) Approve(managerID uuid.UUID) error {
	l.Status = "approved"
	l.ApprovedBy = &managerID
	now := time.Now()
	l.ApprovedDate = &now
	return nil
}

func (l *Leave) Reject(managerID uuid.UUID) {
	l.Status = "rejected"
	l.ApprovedBy = &managerID
	now := time.Now()
	l.ApprovedDate = &now
}

type LeaveBalance struct {
	ID          uuid.UUID `json:"ID"`
	TenantID    uuid.UUID `json:"TenantID"`
	EmployeeID  uuid.UUID `json:"EmployeeID"`
	LeaveTypeID uuid.UUID `json:"LeaveTypeID"`
	Year        int       `json:"Year"`
	TotalDays   int       `json:"TotalDays"`
	UsedDays    int       `json:"UsedDays"`
	PendingDays int       `json:"PendingDays"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
}
