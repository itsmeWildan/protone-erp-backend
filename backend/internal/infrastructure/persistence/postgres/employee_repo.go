package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/employee"
)

// employeeRepo mengimplementasikan domain/employee.Repository (write port)
// dan domain/employee.QueryRepository (read port)
type employeeRepo struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) (employee.Repository, employee.QueryRepository) {
	r := &employeeRepo{pool: pool}
	return r, r
}

// ─── Write Operations ───────────────────────────────────────────────────────

func (r *employeeRepo) Save(ctx context.Context, e *employee.Employee) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		INSERT INTO employees (
			id, tenant_id, user_id, nik, full_name, email, phone,
			gender, birth_date, department_id, position_id, manager_id,
			employment_type, status, join_date, end_date,
			basic_salary, bank_name, bank_account_no, bank_account_name,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12,
			$13, $14, $15, $16,
			$17, $18, $19, $20,
			$21, $22
		)
	`

	_, err := db.Exec(ctx, query,
		e.ID, e.TenantID, e.UserID, e.NIK, e.FullName, e.Email, e.Phone,
		e.Gender, e.BirthDate, e.DepartmentID, e.PositionID, e.ManagerID,
		e.EmploymentType, e.Status, e.JoinDate, e.EndDate,
		e.BasicSalary, e.BankName, e.BankAccountNo, e.BankAccountName,
		e.CreatedAt, e.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("employee.Save: %w", err)
	}
	return nil
}

func (r *employeeRepo) Update(ctx context.Context, e *employee.Employee) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		UPDATE employees SET
			nik = $1, full_name = $2, email = $3, phone = $4,
			gender = $5, birth_date = $6,
			department_id = $7, position_id = $8, manager_id = $9,
			employment_type = $10, status = $11,
			end_date = $12, basic_salary = $13,
			bank_name = $14, bank_account_no = $15, bank_account_name = $16,
			updated_at = $17
		WHERE id = $18 AND tenant_id = $19 AND deleted_at IS NULL
	`

	result, err := db.Exec(ctx, query,
		e.NIK, e.FullName, e.Email, e.Phone,
		e.Gender, e.BirthDate,
		e.DepartmentID, e.PositionID, e.ManagerID,
		e.EmploymentType, e.Status,
		e.EndDate, e.BasicSalary,
		e.BankName, e.BankAccountNo, e.BankAccountName,
		e.UpdatedAt,
		e.ID, e.TenantID,
	)
	if err != nil {
		return fmt.Errorf("employee.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return employee.ErrNotFound
	}
	return nil
}

func (r *employeeRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		UPDATE employees SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`
	result, err := db.Exec(ctx, query, id, tenantID)
	if err != nil {
		return fmt.Errorf("employee.Delete: %w", err)
	}
	if result.RowsAffected() == 0 {
		return employee.ErrNotFound
	}
	return nil
}

func (r *employeeRepo) FindByID(ctx context.Context, tenantID, id uuid.UUID) (*employee.Employee, error) {
	query := `
		SELECT id, tenant_id, user_id, nik, full_name, email, phone,
			gender, birth_date, department_id, position_id, manager_id,
			employment_type, status, join_date, end_date,
			basic_salary, bank_name, bank_account_no, bank_account_name,
			created_at, updated_at
		FROM employees
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
	`

	var e employee.Employee
	var phone, gender, bankName, bankAccNo, bankAccName *string

	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&e.ID, &e.TenantID, &e.UserID, &e.NIK, &e.FullName, &e.Email, &phone,
		&gender, &e.BirthDate, &e.DepartmentID, &e.PositionID, &e.ManagerID,
		&e.EmploymentType, &e.Status, &e.JoinDate, &e.EndDate,
		&e.BasicSalary, &bankName, &bankAccNo, &bankAccName,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, employee.ErrNotFound
		}
		return nil, fmt.Errorf("employee.FindByID: %w", err)
	}

	if phone != nil { e.Phone = *phone }
	if gender != nil { e.Gender = employee.Gender(*gender) }
	if bankName != nil { e.BankName = *bankName }
	if bankAccNo != nil { e.BankAccountNo = *bankAccNo }
	if bankAccName != nil { e.BankAccountName = *bankAccName }

	return &e, nil
}

func (r *employeeRepo) FindByNIK(ctx context.Context, tenantID uuid.UUID, nik string) (*employee.Employee, error) {
	query := `
		SELECT id FROM employees
		WHERE tenant_id = $1 AND nik = $2 AND deleted_at IS NULL
		LIMIT 1
	`
	var e employee.Employee
	err := r.pool.QueryRow(ctx, query, tenantID, nik).Scan(&e.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // tidak ada — boleh digunakan
		}
		return nil, fmt.Errorf("employee.FindByNIK: %w", err)
	}
	return &e, nil
}

func (r *employeeRepo) FindByUserID(ctx context.Context, tenantID, userID uuid.UUID) (*employee.Employee, error) {
	query := `
		SELECT id, tenant_id, user_id, nik, full_name, email, phone,
			gender, birth_date, department_id, position_id, manager_id,
			employment_type, status, join_date, end_date,
			basic_salary, bank_name, bank_account_no, bank_account_name,
			created_at, updated_at
		FROM employees
		WHERE tenant_id = $1 AND user_id = $2 AND deleted_at IS NULL
		LIMIT 1
	`

	var e employee.Employee
	var phone, gender, bankName, bankAccNo, bankAccName *string

	err := r.pool.QueryRow(ctx, query, tenantID, userID).Scan(
		&e.ID, &e.TenantID, &e.UserID, &e.NIK, &e.FullName, &e.Email, &phone,
		&gender, &e.BirthDate, &e.DepartmentID, &e.PositionID, &e.ManagerID,
		&e.EmploymentType, &e.Status, &e.JoinDate, &e.EndDate,
		&e.BasicSalary, &bankName, &bankAccNo, &bankAccName,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("employee.FindByUserID: %w", err)
	}

	if phone != nil { e.Phone = *phone }
	if gender != nil { e.Gender = employee.Gender(*gender) }
	if bankName != nil { e.BankName = *bankName }
	if bankAccNo != nil { e.BankAccountNo = *bankAccNo }
	if bankAccName != nil { e.BankAccountName = *bankAccName }

	return &e, nil
}

// ─── Read Operations (QueryRepository) ─────────────────────────────────────

func (r *employeeRepo) List(ctx context.Context, filter employee.Filter) ([]employee.EmployeeListItem, int64, error) {
	baseWhere := `WHERE e.tenant_id = $1 AND e.deleted_at IS NULL`
	args := []any{filter.TenantID}
	argIdx := 2

	if filter.DepartmentID != nil {
		baseWhere += fmt.Sprintf(" AND e.department_id = $%d", argIdx)
		args = append(args, *filter.DepartmentID)
		argIdx++
	}
	if filter.Status != nil {
		baseWhere += fmt.Sprintf(" AND e.status = $%d", argIdx)
		args = append(args, *filter.Status)
		argIdx++
	}
	if filter.Search != "" {
		baseWhere += fmt.Sprintf(" AND (e.nik ILIKE $%d OR e.full_name ILIKE $%d)", argIdx, argIdx+1)
		search := "%" + filter.Search + "%"
		args = append(args, search, search)
		argIdx += 2
	}

	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM employees e ` + baseWhere
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("employee.List count: %w", err)
	}

	// Data query dengan pagination
	offset := (filter.Page - 1) * filter.PerPage
	dataArgs := append(args, filter.PerPage, offset)
	dataQuery := fmt.Sprintf(`
		SELECT e.id, e.nik, e.full_name, e.email,
			d.name AS department_name,
			p.name AS position_name,
			e.status, e.join_date::TEXT, e.employment_type
		FROM employees e
		JOIN departments d ON d.id = e.department_id
		JOIN positions p ON p.id = e.position_id
		%s
		ORDER BY e.full_name ASC
		LIMIT $%d OFFSET $%d
	`, baseWhere, argIdx, argIdx+1)

	rows, err := r.pool.Query(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("employee.List query: %w", err)
	}
	defer rows.Close()

	var result []employee.EmployeeListItem
	for rows.Next() {
		var item employee.EmployeeListItem
		if err := rows.Scan(
			&item.ID, &item.NIK, &item.FullName, &item.Email,
			&item.DepartmentName, &item.PositionName,
			&item.Status, &item.JoinDate, &item.EmploymentType,
		); err != nil {
			return nil, 0, fmt.Errorf("employee.List scan: %w", err)
		}
		result = append(result, item)
	}

	return result, total, nil
}

func (r *employeeRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*employee.EmployeeDetail, error) {
	query := `
		SELECT e.id, e.nik, e.full_name, e.email, COALESCE(e.phone, ''), COALESCE(e.gender::TEXT, ''),
			e.birth_date::TEXT,
			e.department_id, d.name,
			e.position_id, p.name,
			e.manager_id, m.full_name,
			e.employment_type, e.status,
			e.join_date::TEXT,
			COALESCE(e.basic_salary, 0), COALESCE(e.bank_name, ''), COALESCE(e.bank_account_no, ''), COALESCE(e.bank_account_name, '')
		FROM employees e
		JOIN departments d ON d.id = e.department_id
		JOIN positions p ON p.id = e.position_id
		LEFT JOIN employees m ON m.id = e.manager_id
		WHERE e.id = $1 AND e.tenant_id = $2 AND e.deleted_at IS NULL
	`

	var det employee.EmployeeDetail
	err := r.pool.QueryRow(ctx, query, id, tenantID).Scan(
		&det.ID, &det.NIK, &det.FullName, &det.Email, &det.Phone, &det.Gender,
		&det.BirthDate,
		&det.DepartmentID, &det.DepartmentName,
		&det.PositionID, &det.PositionName,
		&det.ManagerID, &det.ManagerName,
		&det.EmploymentType, &det.Status,
		&det.JoinDate,
		&det.BasicSalary, &det.BankName, &det.BankAccountNo, &det.BankAccountName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, employee.ErrNotFound
		}
		return nil, fmt.Errorf("employee.GetByID: %w", err)
	}
	return &det, nil
}
