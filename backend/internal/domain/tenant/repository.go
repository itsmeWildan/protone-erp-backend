package tenant

import (
	"context"

	"github.com/google/uuid"
)

// Repository adalah port untuk persistence tenant (Write).
type Repository interface {
	Create(ctx context.Context, t *Tenant) error
	Update(ctx context.Context, t *Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetBySlug(ctx context.Context, slug string) (*Tenant, error)
}

// QueryRepository adalah port untuk read-only operations (semi-CQRS style).
type QueryRepository interface {
	List(ctx context.Context) ([]*Tenant, error)
}
