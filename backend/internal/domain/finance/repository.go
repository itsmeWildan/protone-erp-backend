package finance

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// COA
	GetCOAByCode(ctx context.Context, tenantID uuid.UUID, code string) (*ChartOfAccount, error)
	
	// Journal
	CreateJournal(ctx context.Context, j *JournalEntry) error
	GetJournalBySource(ctx context.Context, sourceType string, sourceID uuid.UUID) (*JournalEntry, error)

	// Budgeting
	GetBudget(ctx context.Context, tenantID, deptID uuid.UUID, month, year int) (*DepartmentBudget, error)
	UpdateBudget(ctx context.Context, b *DepartmentBudget) error
}
