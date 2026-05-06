package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/protone/erp/internal/domain/user"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) user.Repository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, u *user.User) error {
	db := ExtractTx(ctx, r.pool)

	query := `
		INSERT INTO users (id, tenant_id, name, email, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(ctx, query, u.ID, u.TenantID, u.Name, u.Email, u.PasswordHash, u.Role, u.CreatedAt)
	if err != nil {
		return fmt.Errorf("user.Create: %w", err)
	}
	return nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, tenant_id, name, email, password_hash, role, created_at FROM users WHERE email = $1`

	var u user.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Not found is okay for registration check
		}
		return nil, fmt.Errorf("user.GetByEmail: %w", err)
	}
	return &u, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	query := `SELECT id, tenant_id, name, email, password_hash, role, created_at FROM users WHERE id = $1`

	var u user.User
	err := r.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.TenantID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return &u, nil
}
