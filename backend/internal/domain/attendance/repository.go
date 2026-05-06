package attendance

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, a *Attendance) error
	Update(ctx context.Context, a *Attendance) error
	GetByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*Attendance, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Attendance, error)
}

type QueryRepository interface {
	GetByEmployee(ctx context.Context, tenantID, employeeID uuid.UUID, startDate, endDate time.Time) ([]Attendance, error)
	GetStats(ctx context.Context, tenantID uuid.UUID, date string) (int, int, int, int, error)
}
