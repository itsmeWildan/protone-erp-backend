package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/leave"
)

type leaveRepo struct {
	pool *pgxpool.Pool
}

func NewLeaveRepository(pool *pgxpool.Pool) leave.Repository {
	return &leaveRepo{pool: pool}
}

func (r *leaveRepo) GetTypes(ctx context.Context, tenantID uuid.UUID) ([]leave.LeaveType, error) {
	query := `SELECT id, name, description, max_days FROM leave_types WHERE tenant_id = $1`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil { return nil, err }
	defer rows.Close()

	result := []leave.LeaveType{}
	for rows.Next() {
		var lt leave.LeaveType
		if err := rows.Scan(&lt.ID, &lt.Name, &lt.Description, &lt.MaxDays); err != nil { return nil, err }
		result = append(result, lt)
	}
	return result, nil
}

func (r *leaveRepo) Save(ctx context.Context, l *leave.Leave) error {
	query := `INSERT INTO leaves (id, tenant_id, employee_id, leave_type_id, start_date, end_date, total_days, reason, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.pool.Exec(ctx, query, l.ID, l.TenantID, l.EmployeeID, l.LeaveTypeID, l.StartDate, l.EndDate, l.TotalDays, l.Reason, l.Status, l.CreatedAt, l.UpdatedAt)
	return err
}

func (r *leaveRepo) GetByID(ctx context.Context, id uuid.UUID) (*leave.Leave, error) {
	query := `SELECT id, tenant_id, employee_id, leave_type_id, start_date, end_date, total_days, reason, status FROM leaves WHERE id = $1`
	var l leave.Leave
	err := r.pool.QueryRow(ctx, query, id).Scan(&l.ID, &l.TenantID, &l.EmployeeID, &l.LeaveTypeID, &l.StartDate, &l.EndDate, &l.TotalDays, &l.Reason, &l.Status)
	if err != nil { return nil, err }
	return &l, nil
}

func (r *leaveRepo) Update(ctx context.Context, l *leave.Leave) error {
	query := `UPDATE leaves SET status = $1, approved_by = $2, approved_at = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, l.Status, l.ApprovedBy, l.ApprovedDate, l.ID)
	return err
}

func (r *leaveRepo) GetBalance(ctx context.Context, employeeID, leaveTypeID uuid.UUID, year int) (*leave.LeaveBalance, error) {
	query := `SELECT id, total_days, used_days, pending_days FROM leave_balances WHERE employee_id = $1 AND leave_type_id = $2 AND year = $3`
	var b leave.LeaveBalance
	err := r.pool.QueryRow(ctx, query, employeeID, leaveTypeID, year).Scan(&b.ID, &b.TotalDays, &b.UsedDays, &b.PendingDays)
	if err != nil { return nil, err }
	b.EmployeeID = employeeID
	b.LeaveTypeID = leaveTypeID
	b.Year = year
	return &b, nil
}

func (r *leaveRepo) CreateBalance(ctx context.Context, b *leave.LeaveBalance) error {
	query := `INSERT INTO leave_balances (id, tenant_id, employee_id, leave_type_id, year, total_days, used_days, pending_days, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.pool.Exec(ctx, query, b.ID, b.TenantID, b.EmployeeID, b.LeaveTypeID, b.Year, b.TotalDays, b.UsedDays, b.PendingDays, b.UpdatedAt)
	return err
}

func (r *leaveRepo) UpdateBalance(ctx context.Context, b *leave.LeaveBalance) error {
	query := `UPDATE leave_balances SET used_days = $1, pending_days = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, b.UsedDays, b.PendingDays, b.ID)
	return err
}

func (r *leaveRepo) GetByEmployee(ctx context.Context, employeeID uuid.UUID) ([]leave.Leave, error) {
	query := `SELECT id, leave_type_id, start_date, end_date, total_days, reason, status FROM leaves WHERE employee_id = $1`
	rows, err := r.pool.Query(ctx, query, employeeID)
	if err != nil { return nil, err }
	defer rows.Close()

	result := []leave.Leave{}
	for rows.Next() {
		var l leave.Leave
		if err := rows.Scan(&l.ID, &l.LeaveTypeID, &l.StartDate, &l.EndDate, &l.TotalDays, &l.Reason, &l.Status); err != nil { return nil, err }
		result = append(result, l)
	}
	return result, nil
}
