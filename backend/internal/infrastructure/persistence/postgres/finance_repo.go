package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/finance"
)

type financeRepo struct {
	pool *pgxpool.Pool
}

func NewFinanceRepository(pool *pgxpool.Pool) finance.Repository {
	return &financeRepo{pool: pool}
}

func (r *financeRepo) GetCOAByCode(ctx context.Context, tenantID uuid.UUID, code string) (*finance.ChartOfAccount, error) {
	query := `SELECT id, tenant_id, code, name, type, normal_balance FROM chart_of_accounts WHERE tenant_id = $1 AND code = $2`
	var coa finance.ChartOfAccount
	err := r.pool.QueryRow(ctx, query, tenantID, code).Scan(&coa.ID, &coa.TenantID, &coa.Code, &coa.Name, &coa.Type, &coa.NormalBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &coa, nil
}

func (r *financeRepo) CreateJournal(ctx context.Context, j *finance.JournalEntry) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Insert Header
	queryHeader := `
		INSERT INTO journal_entries (id, tenant_id, journal_no, date, description, status, source_type, source_id, total_debit, total_credit)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.Exec(ctx, queryHeader, j.ID, j.TenantID, j.JournalNo, j.Date, j.Description, j.Status, j.SourceType, j.SourceID, j.TotalDebit, j.TotalCredit)
	if err != nil {
		return err
	}

	// 2. Insert Lines
	for i, l := range j.Lines {
		queryLine := `
			INSERT INTO journal_lines (id, journal_entry_id, coa_id, description, debit, credit, line_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.Exec(ctx, queryLine, l.ID, j.ID, l.COAID, l.Description, l.Debit, l.Credit, i)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *financeRepo) GetJournalBySource(ctx context.Context, sourceType string, sourceID uuid.UUID) (*finance.JournalEntry, error) {
	query := `SELECT id, tenant_id, journal_no, date, description, status, total_debit, total_credit FROM journal_entries WHERE source_type = $1 AND source_id = $2`
	var j finance.JournalEntry
	err := r.pool.QueryRow(ctx, query, sourceType, sourceID).Scan(&j.ID, &j.TenantID, &j.JournalNo, &j.Date, &j.Description, &j.Status, &j.TotalDebit, &j.TotalCredit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

// Budgeting Implementation
func (r *financeRepo) GetBudget(ctx context.Context, tenantID, deptID uuid.UUID, month, year int) (*finance.DepartmentBudget, error) {
	query := `SELECT id, tenant_id, department_id, month, year, allocated_amount, spent_amount, created_at, updated_at FROM department_budgets WHERE tenant_id = $1 AND department_id = $2 AND month = $3 AND year = $4`
	var b finance.DepartmentBudget
	err := r.pool.QueryRow(ctx, query, tenantID, deptID, month, year).Scan(&b.ID, &b.TenantID, &b.DepartmentID, &b.Month, &b.Year, &b.AllocatedAmount, &b.SpentAmount, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

func (r *financeRepo) UpdateBudget(ctx context.Context, b *finance.DepartmentBudget) error {
	query := `
		INSERT INTO department_budgets (id, tenant_id, department_id, month, year, allocated_amount, spent_amount, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (department_id, month, year) DO UPDATE SET
			allocated_amount = EXCLUDED.allocated_amount,
			spent_amount = EXCLUDED.spent_amount,
			updated_at = NOW()
	`
	_, err := r.pool.Exec(ctx, query, b.ID, b.TenantID, b.DepartmentID, b.Month, b.Year, b.AllocatedAmount, b.SpentAmount)
	return err
}
