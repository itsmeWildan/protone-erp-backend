package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/tenant"
)

type tenantRepo struct {
	pool *pgxpool.Pool
}

func NewTenantRepository(pool *pgxpool.Pool) tenant.Repository {
	return &tenantRepo{pool: pool}
}

func (r *tenantRepo) Create(ctx context.Context, t *tenant.Tenant) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		INSERT INTO tenants (id, name, slug, email, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(ctx, query, t.ID, t.Name, t.Slug, t.Email, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("tenant.Create: %w", err)
	}
	return nil
}

func (r *tenantRepo) Update(ctx context.Context, t *tenant.Tenant) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		UPDATE tenants SET name = $1, status = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := db.Exec(ctx, query, t.Name, t.Status, t.UpdatedAt, t.ID)
	if err != nil {
		return fmt.Errorf("tenant.Update: %w", err)
	}
	if result.RowsAffected() == 0 {
		return tenant.ErrTenantNotFound
	}
	return nil
}

func (r *tenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*tenant.Tenant, error) {
	query := `SELECT id, name, slug, email, status, created_at, updated_at FROM tenants WHERE id = $1`

	var t tenant.Tenant
	err := r.pool.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Slug, &t.Email, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, tenant.ErrTenantNotFound
		}
		return nil, fmt.Errorf("tenant.GetByID: %w", err)
	}
	return &t, nil
}

func (r *tenantRepo) GetBySlug(ctx context.Context, slug string) (*tenant.Tenant, error) {
	query := `SELECT id, name, slug, email, status, created_at, updated_at FROM tenants WHERE slug = $1`

	var t tenant.Tenant
	err := r.pool.QueryRow(ctx, query, slug).Scan(&t.ID, &t.Name, &t.Slug, &t.Email, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Tidak ada slug yang sama — boleh digunakan
		}
		return nil, fmt.Errorf("tenant.GetBySlug: %w", err)
	}
	return &t, nil
}
