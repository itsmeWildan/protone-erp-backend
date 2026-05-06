package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/reimbursement"
)

type reimbursementRepo struct {
	pool *pgxpool.Pool
}

func NewReimbursementRepository(pool *pgxpool.Pool) reimbursement.Repository {
	return &reimbursementRepo{pool: pool}
}

func (r *reimbursementRepo) Save(ctx context.Context, claim *reimbursement.Reimbursement) error {
	query := `INSERT INTO reimbursements (id, tenant_id, employee_id, date, category, amount, description, receipt_url, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, claim.ID, claim.TenantID, claim.EmployeeID, claim.Date, claim.Category, claim.Amount, claim.Description, claim.ReceiptURL, claim.Status, claim.CreatedAt, claim.UpdatedAt)
	return err
}

func (r *reimbursementRepo) GetByID(ctx context.Context, id uuid.UUID) (*reimbursement.Reimbursement, error) {
	query := `SELECT id, tenant_id, employee_id, date, category, amount, description, receipt_url, status, approved_by, approved_at, created_at, updated_at FROM reimbursements WHERE id = $1`
	var c reimbursement.Reimbursement
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.TenantID, &c.EmployeeID, &c.Date, &c.Category, &c.Amount, &c.Description, &c.ReceiptURL, &c.Status, 
		&c.ApprovedBy, &c.ApprovedAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil { return nil, err }
	return &c, nil
}

func (r *reimbursementRepo) Update(ctx context.Context, claim *reimbursement.Reimbursement) error {
	query := `UPDATE reimbursements SET status = $1, approved_by = $2, approved_at = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, claim.Status, claim.ApprovedBy, claim.ApprovedAt, claim.ID)
	return err
}

func (r *reimbursementRepo) GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]reimbursement.Reimbursement, error) {
	query := `SELECT id, tenant_id, employee_id, date, category, amount, description, receipt_url, status, approved_by, approved_at, created_at, updated_at FROM reimbursements WHERE employee_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	result := []reimbursement.Reimbursement{}
	for rows.Next() {
		var c reimbursement.Reimbursement
		err := rows.Scan(
			&c.ID, &c.TenantID, &c.EmployeeID, &c.Date, &c.Category, &c.Amount, &c.Description, &c.ReceiptURL, &c.Status,
			&c.ApprovedBy, &c.ApprovedAt, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil { return nil, err }
		result = append(result, c)
	}
	return result, nil
}
