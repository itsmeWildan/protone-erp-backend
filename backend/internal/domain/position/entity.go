package position

import (
	"time"

	"github.com/google/uuid"
)

type Position struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	Name      string
	Level     int // Misal: 1 untuk Staff, 10 untuk Director
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(tenantID uuid.UUID, name string, level int) *Position {
	now := time.Now()
	return &Position{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      name,
		Level:     level,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
