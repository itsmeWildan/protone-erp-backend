package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/attendance"
)

type attendanceRepo struct {
	pool *pgxpool.Pool
}

func NewAttendanceRepository(pool *pgxpool.Pool) (attendance.Repository, attendance.QueryRepository) {
	r := &attendanceRepo{pool: pool}
	return r, r
}

func (r *attendanceRepo) Save(ctx context.Context, a *attendance.Attendance) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		INSERT INTO attendances (
			id, tenant_id, employee_id, date, check_in, status, location_in, notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := db.Exec(ctx, query,
		a.ID, a.TenantID, a.EmployeeID, a.Date, a.CheckIn, a.Status, a.LocationIn, a.Notes, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("attendance.Save: %w", err)
	}
	return nil
}

func (r *attendanceRepo) Update(ctx context.Context, a *attendance.Attendance) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		UPDATE attendances SET
			check_out = $1, location_out = $2, status = $3, updated_at = $4
		WHERE id = $5
	`

	_, err := db.Exec(ctx, query,
		a.CheckOut, a.LocationOut, a.Status, a.UpdatedAt, a.ID,
	)
	if err != nil {
		return fmt.Errorf("attendance.Update: %w", err)
	}
	return nil
}

func (r *attendanceRepo) GetByEmployeeAndDate(ctx context.Context, employeeID uuid.UUID, date time.Time) (*attendance.Attendance, error) {
	query := `
		SELECT id, tenant_id, employee_id, date, check_in, check_out, status, location_in, location_out, notes, created_at, updated_at
		FROM attendances
		WHERE employee_id = $1 AND date = $2
	`

	var a attendance.Attendance
	err := r.pool.QueryRow(ctx, query, employeeID, date.Format("2006-01-02")).Scan(
		&a.ID, &a.TenantID, &a.EmployeeID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Status, &a.LocationIn, &a.LocationOut, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("attendance.GetByEmployeeAndDate: %w", err)
	}
	return &a, nil
}

func (r *attendanceRepo) GetByID(ctx context.Context, id uuid.UUID) (*attendance.Attendance, error) {
	query := `
		SELECT id, tenant_id, employee_id, date, check_in, check_out, status, location_in, location_out, notes, created_at, updated_at
		FROM attendances
		WHERE id = $1
	`

	var a attendance.Attendance
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.TenantID, &a.EmployeeID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Status, &a.LocationIn, &a.LocationOut, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("attendance not found")
		}
		return nil, fmt.Errorf("attendance.GetByID: %w", err)
	}
	return &a, nil
}

func (r *attendanceRepo) GetByEmployee(ctx context.Context, tenantID, employeeID uuid.UUID, startDate, endDate time.Time) ([]attendance.Attendance, error) {
	query := `
		SELECT id, tenant_id, employee_id, date, check_in, check_out, status, location_in, location_out, notes, created_at, updated_at
		FROM attendances
		WHERE tenant_id = $1 AND employee_id = $2 AND date BETWEEN $3 AND $4
		ORDER BY date DESC
	`

	rows, err := r.pool.Query(ctx, query, tenantID, employeeID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("attendance.GetByEmployee: %w", err)
	}
	defer rows.Close()

	var result []attendance.Attendance
	for rows.Next() {
		var a attendance.Attendance
		err := rows.Scan(
			&a.ID, &a.TenantID, &a.EmployeeID, &a.Date, &a.CheckIn, &a.CheckOut, &a.Status, &a.LocationIn, &a.LocationOut, &a.Notes, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("attendance.GetByEmployee scan: %w", err)
		}
		result = append(result, a)
	}
	return result, nil
}

func (r *attendanceRepo) GetStats(ctx context.Context, tenantID uuid.UUID, date string) (int, int, int, int, error) {
	var total, present, late, onLeave int

	// 1. Total Employees
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM employees WHERE tenant_id = $1", tenantID).Scan(&total)

	// 2. Present Today
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM attendances WHERE tenant_id = $1 AND date = $2", tenantID, date).Scan(&present)

	// 3. Late Today
	_ = r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM attendances WHERE tenant_id = $1 AND date = $2 AND check_in > '09:00:00'", tenantID, date).Scan(&late)

	// 4. On Leave Today
	queryLeave := `
		SELECT COUNT(*) FROM leaves 
		WHERE tenant_id = $1 AND status = 'approved' 
		AND $2 BETWEEN start_date AND end_date
	`
	_ = r.pool.QueryRow(ctx, queryLeave, tenantID, date).Scan(&onLeave)

	return total, present, late, onLeave, nil
}
