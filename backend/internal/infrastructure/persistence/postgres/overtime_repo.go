package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/overtime"
)

type overtimeRepo struct {
	pool *pgxpool.Pool
}

func NewOvertimeRepository(pool *pgxpool.Pool) overtime.Repository {
	return &overtimeRepo{pool: pool}
}

func (r *overtimeRepo) Save(ctx context.Context, o *overtime.OvertimeRequest) error {
	query := `
		INSERT INTO overtime_requests (
			id, tenant_id, employee_id, date, start_time, end_time, duration_hours, reason, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.pool.Exec(ctx, query, o.ID, o.TenantID, o.EmployeeID, o.Date, o.StartTime, o.EndTime, o.DurationHours, o.Reason, o.Status, o.CreatedAt, o.UpdatedAt)
	return err
}

func (r *overtimeRepo) Update(ctx context.Context, o *overtime.OvertimeRequest) error {
	query := `
		UPDATE overtime_requests SET status = $1, approved_by = $2, approved_at = $3, rejection_note = $4, updated_at = NOW()
		WHERE id = $5
	`
	_, err := r.pool.Exec(ctx, query, o.Status, o.ApprovedBy, o.ApprovedAt, o.RejectionNote, o.ID)
	return err
}

func (r *overtimeRepo) GetByID(ctx context.Context, id uuid.UUID) (*overtime.OvertimeRequest, error) {
	query := `SELECT id, tenant_id, employee_id, date, start_time, end_time, duration_hours, reason, status, approved_by, approved_at, rejection_note, created_at, updated_at FROM overtime_requests WHERE id = $1`
	var o overtime.OvertimeRequest
	var startTime, endTime time.Time // pgx maps TIME to time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(&o.ID, &o.TenantID, &o.EmployeeID, &o.Date, &startTime, &endTime, &o.DurationHours, &o.Reason, &o.Status, &o.ApprovedBy, &o.ApprovedAt, &o.RejectionNote, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	o.StartTime = startTime.Format("15:04")
	o.EndTime = endTime.Format("15:04")
	return &o, nil
}

func (r *overtimeRepo) GetByEmployeeAndMonth(ctx context.Context, employeeID uuid.UUID, month, year int) ([]overtime.OvertimeRequest, error) {
	query := `
		SELECT id, tenant_id, employee_id, date, start_time, end_time, duration_hours, reason, status, approved_by, approved_at, rejection_note, created_at, updated_at 
		FROM overtime_requests 
		WHERE employee_id = $1 AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(YEAR FROM date) = $3
		ORDER BY date DESC
	`
	rows, err := r.pool.Query(ctx, query, employeeID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []overtime.OvertimeRequest{}
	for rows.Next() {
		var o overtime.OvertimeRequest
		var st, et time.Time
		if err := rows.Scan(&o.ID, &o.TenantID, &o.EmployeeID, &o.Date, &st, &et, &o.DurationHours, &o.Reason, &o.Status, &o.ApprovedBy, &o.ApprovedAt, &o.RejectionNote, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.StartTime = st.Format("15:04")
		o.EndTime = et.Format("15:04")
		list = append(list, o)
	}
	return list, nil
}

func (r *overtimeRepo) GetApprovedSumByEmployee(ctx context.Context, employeeID uuid.UUID, month, year int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(duration_hours), 0) 
		FROM overtime_requests 
		WHERE employee_id = $1 AND status = 'approved' 
		AND EXTRACT(MONTH FROM date) = $2 AND EXTRACT(YEAR FROM date) = $3
	`
	var total float64
	err := r.pool.QueryRow(ctx, query, employeeID, month, year).Scan(&total)
	return total, err
}
