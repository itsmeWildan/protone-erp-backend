package reimbursement

import (
	"time"

	"github.com/google/uuid"
)

type Reimbursement struct {
	ID          uuid.UUID  `json:"ID"`
	TenantID    uuid.UUID  `json:"TenantID"`
	EmployeeID  uuid.UUID  `json:"EmployeeID"`
	Date        time.Time  `json:"Date"`
	Category    string     `json:"Category"`
	Amount      float64    `json:"Amount"`
	Description string     `json:"Description"`
	ReceiptURL  string     `json:"ReceiptURL"`
	Status      string     `json:"Status"`
	ApprovedBy  *uuid.UUID `json:"ApprovedBy"`
	ApprovedAt  *time.Time `json:"ApprovedAt"`
	CreatedAt   time.Time  `json:"CreatedAt"`
	UpdatedAt   time.Time  `json:"UpdatedAt"`
}

func NewReimbursement(tenantID, empID uuid.UUID, category string, amount float64, date time.Time, desc string) *Reimbursement {
	return &Reimbursement{
		ID:          uuid.New(),
		TenantID:    tenantID,
		EmployeeID:  empID,
		Date:        date,
		Category:    category,
		Amount:      amount,
		Description: desc,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (r *Reimbursement) Approve(managerID uuid.UUID) error {
	r.Status = "approved"
	r.ApprovedBy = &managerID
	now := time.Now()
	r.ApprovedAt = &now
	return nil
}
