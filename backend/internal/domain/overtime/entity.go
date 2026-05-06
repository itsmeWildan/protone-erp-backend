package overtime

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type OvertimeRequest struct {
	ID            uuid.UUID  `json:"ID"`
	TenantID      uuid.UUID  `json:"TenantID"`
	EmployeeID    uuid.UUID  `json:"EmployeeID"`
	Date          time.Time  `json:"Date"`
	StartTime     string     `json:"StartTime"`
	EndTime       string     `json:"EndTime"`
	DurationHours float64    `json:"DurationHours"`
	Reason        string     `json:"Reason"`
	Status        Status     `json:"Status"`
	ApprovedBy    *uuid.UUID `json:"ApprovedBy"`
	ApprovedAt    *time.Time `json:"ApprovedAt"`
	RejectionNote *string    `json:"RejectionNote"`
	CreatedAt     time.Time  `json:"CreatedAt"`
	UpdatedAt     time.Time  `json:"UpdatedAt"`
}

func NewOvertimeRequest(tenantID, employeeID uuid.UUID, date time.Time, start, end string, duration float64, reason string) *OvertimeRequest {
	now := time.Now()
	return &OvertimeRequest{
		ID:            uuid.New(),
		TenantID:      tenantID,
		EmployeeID:    employeeID,
		Date:          date,
		StartTime:     start,
		EndTime:       end,
		DurationHours: duration,
		Reason:        reason,
		Status:        StatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (o *OvertimeRequest) Approve(managerUserID uuid.UUID) {
	now := time.Now()
	o.Status = StatusApproved
	o.ApprovedBy = &managerUserID
	o.ApprovedAt = &now
	o.UpdatedAt = now
}

func (o *OvertimeRequest) Reject(note string) {
	o.Status = StatusRejected
	o.RejectionNote = &note
	o.UpdatedAt = time.Now()
}
