package department

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(tenantID uuid.UUID, name, code string) *Department {
	now := time.Now()
	return &Department{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		Code:      code,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
