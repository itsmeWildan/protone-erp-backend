package finance

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type DepartmentBudget struct {
	ID              uuid.UUID
	TenantID        uuid.UUID
	DepartmentID    uuid.UUID
	Month           int
	Year            int
	AllocatedAmount float64
	SpentAmount     float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (b *DepartmentBudget) RemainingBalance() float64 {
	return b.AllocatedAmount - b.SpentAmount
}

func (b *DepartmentBudget) Deduct(amount float64) error {
	if b.RemainingBalance() < amount {
		return fmt.Errorf("insufficient budget for department. Remaining: %.2f, Requested: %.2f", b.RemainingBalance(), amount)
	}
	b.SpentAmount += amount
	b.UpdatedAt = time.Now()
	return nil
}
