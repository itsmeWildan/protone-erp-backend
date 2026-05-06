package attendance

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPresent    Status = "present"
	StatusAbsent     Status = "absent"
	StatusLate       Status = "late"
	StatusEarlyLeave Status = "early_leave"
	StatusOnLeave    Status = "on_leave"
)

var (
	ErrAlreadyClockedIn  = errors.New("already clocked in for today")
	ErrNotClockedIn      = errors.New("not clocked in yet")
	ErrAlreadyClockedOut = errors.New("already clocked out for today")
)

type Attendance struct {
	ID          uuid.UUID  `json:"ID"`
	TenantID    uuid.UUID  `json:"TenantID"`
	EmployeeID  uuid.UUID  `json:"EmployeeID"`
	Date        time.Time  `json:"Date"`
	CheckIn     *time.Time `json:"CheckIn"`
	CheckOut    *time.Time `json:"CheckOut"`
	Status      Status     `json:"Status"`
	LocationIn  *string    `json:"LocationIn"`
	LocationOut *string    `json:"LocationOut"`
	Notes       *string    `json:"Notes"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
}

func NewAttendance(tenantID, employeeID uuid.UUID, locationIn, notes string) *Attendance {
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	locIn := &locationIn
	if locationIn == "" {
		locIn = nil
	}
	n := &notes
	if notes == "" {
		n = nil
	}

	return &Attendance{
		ID:         uuid.New(),
		TenantID:   tenantID,
		EmployeeID: employeeID,
		Date:       date,
		CheckIn:    &now,
		Status:     StatusPresent,
		LocationIn: locIn,
		Notes:      n,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (a *Attendance) ClockOut(locationOut string) error {
	if a.CheckOut != nil {
		return ErrAlreadyClockedOut
	}
	now := time.Now()
	a.CheckOut = &now
	a.LocationOut = &locationOut
	a.UpdatedAt = now
	return nil
}
