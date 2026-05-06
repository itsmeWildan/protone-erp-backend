package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/payroll"
)

type payrollRepo struct {
	pool *pgxpool.Pool
}

func NewPayrollRepository(pool *pgxpool.Pool) payroll.Repository {
	return &payrollRepo{pool: pool}
}

func (r *payrollRepo) GetComponents(ctx context.Context, tenantID uuid.UUID) ([]payroll.SalaryComponent, error) {
	query := `SELECT id, name, type, code, default_amount, is_taxable FROM salary_components WHERE tenant_id = $1`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []payroll.SalaryComponent{}
	for rows.Next() {
		var c payroll.SalaryComponent
		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.Code, &c.DefaultAmount, &c.IsTaxable); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *payrollRepo) GetPeriodByID(ctx context.Context, id uuid.UUID) (*payroll.PayrollPeriod, error) {
	query := `SELECT id, period_month, period_year, status, total_amount FROM payroll_periods WHERE id = $1`
	var p payroll.PayrollPeriod
	err := r.pool.QueryRow(ctx, query, id).Scan(&p.ID, &p.Month, &p.Year, &p.Status, &p.TotalAmount)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

func (r *payrollRepo) GetPeriod(ctx context.Context, tenantID uuid.UUID, month, year int) (*payroll.PayrollPeriod, error) {
	query := `SELECT id, period_month, period_year, status, total_amount FROM payroll_periods WHERE tenant_id = $1 AND period_month = $2 AND period_year = $3`
	var p payroll.PayrollPeriod
	err := r.pool.QueryRow(ctx, query, tenantID, month, year).Scan(&p.ID, &p.Month, &p.Year, &p.Status, &p.TotalAmount)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	p.TenantID = tenantID
	return &p, nil
}

func (r *payrollRepo) CreatePeriod(ctx context.Context, p *payroll.PayrollPeriod) error {
	query := `INSERT INTO payroll_periods (id, tenant_id, period_month, period_year, status, total_amount, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.pool.Exec(ctx, query, p.ID, p.TenantID, p.Month, p.Year, p.Status, p.TotalAmount, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *payrollRepo) UpdatePeriod(ctx context.Context, p *payroll.PayrollPeriod) error {
	query := `UPDATE payroll_periods SET status = $1, total_amount = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, p.Status, p.TotalAmount, p.ID)
	return err
}

func (r *payrollRepo) UpdatePeriodStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE payroll_periods SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func (r *payrollRepo) GetSlipsByPeriod(ctx context.Context, periodID uuid.UUID) ([]payroll.PayrollSlip, error) {
	return []payroll.PayrollSlip{}, nil
}

func (r *payrollRepo) GetSlipByEmployee(ctx context.Context, periodID, employeeID uuid.UUID) (*payroll.PayrollSlip, error) {
	query := `SELECT id, tenant_id, payroll_period_id, employee_id, basic_salary, total_allowance, total_deduction, net_salary, working_days, present_days, overtime_hours, overtime_amount, created_at, updated_at 
	          FROM payroll_slips WHERE payroll_period_id = $1 AND employee_id = $2`
	var s payroll.PayrollSlip
	err := r.pool.QueryRow(ctx, query, periodID, employeeID).Scan(
		&s.ID, &s.TenantID, &s.PayrollPeriodID, &s.EmployeeID, &s.BasicSalary, &s.TotalAllowance, &s.TotalDeduction, &s.NetSalary, &s.WorkingDays, &s.PresentDays, &s.OvertimeHours, &s.OvertimeAmount, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	detQuery := `SELECT id, salary_component_id, type, amount FROM payroll_slip_details WHERE payroll_slip_id = $1`
	rows, _ := r.pool.Query(ctx, detQuery, s.ID)
	defer rows.Close()
	s.Details = []payroll.PayrollSlipDetail{}
	for rows.Next() {
		var d payroll.PayrollSlipDetail
		rows.Scan(&d.ID, &d.SalaryComponentID, &d.Type, &d.Amount)
		s.Details = append(s.Details, d)
	}
	return &s, nil
}

func (r *payrollRepo) SaveSlip(ctx context.Context, s *payroll.PayrollSlip) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	q1 := `INSERT INTO payroll_slips (id, tenant_id, payroll_period_id, employee_id, basic_salary, total_allowance, total_deduction, net_salary, working_days, present_days, overtime_hours, overtime_amount, created_at, updated_at) 
	       VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	_, err = tx.Exec(ctx, q1, s.ID, s.TenantID, s.PayrollPeriodID, s.EmployeeID, s.BasicSalary, s.TotalAllowance, s.TotalDeduction, s.NetSalary, s.WorkingDays, s.PresentDays, s.OvertimeHours, s.OvertimeAmount, s.CreatedAt, s.UpdatedAt)
	if err != nil { return err }

	for _, d := range s.Details {
		q2 := `INSERT INTO payroll_slip_details (id, payroll_slip_id, salary_component_id, type, amount) VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx, q2, uuid.New(), s.ID, d.SalaryComponentID, d.Type, d.Amount)
		if err != nil { return err }
	}

	return tx.Commit(ctx)
}

func (r *payrollRepo) DeleteSlipsByPeriod(ctx context.Context, periodID uuid.UUID) error {
	query := `DELETE FROM payroll_slips WHERE payroll_period_id = $1`
	_, err := r.pool.Exec(ctx, query, periodID)
	return err
}

func (r *payrollRepo) GetDepartmentBreakdown(ctx context.Context, tenantID uuid.UUID, month, year int) (map[uuid.UUID]float64, error) {
	query := `
		SELECT e.department_id, SUM(ps.net_salary)
		FROM payroll_slips ps
		JOIN employees e ON ps.employee_id = e.id
		JOIN payroll_periods pp ON ps.payroll_period_id = pp.id
		WHERE pp.tenant_id = $1 AND pp.period_month = $2 AND pp.period_year = $3
		GROUP BY e.department_id
	`
	rows, err := r.pool.Query(ctx, query, tenantID, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(map[uuid.UUID]float64)
	for rows.Next() {
		var deptID uuid.UUID
		var total float64
		if err := rows.Scan(&deptID, &total); err != nil {
			return nil, err
		}
		res[deptID] = total
	}
	return res, nil
}
