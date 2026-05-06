package tenant

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusTrial    Status = "trial"
)

// Tenant adalah root entity untuk multi-tenancy.
type Tenant struct {
	ID        uuid.UUID
	Name      string
	Slug      string // Alias unik untuk URL atau identifikasi (misal: "acme-corp")
	Email     string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateParams struct {
	Name  string
	Slug  string
	Email string
}

// NewTenant adalah factory untuk membuat tenant baru.
func NewTenant(p CreateParams) (*Tenant, error) {
	if p.Name == "" {
		return nil, ErrNameRequired
	}
	if p.Slug == "" {
		return nil, ErrSlugRequired
	}

	now := time.Now()
	return &Tenant{
		ID:        uuid.New(),
		Name:      p.Name,
		Slug:      p.Slug,
		Email:     p.Email,
		Status:    StatusTrial, // Default baru daftar adalah trial
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (t *Tenant) Activate() {
	t.Status = StatusActive
	t.UpdatedAt = time.Now()
}

func (t *Tenant) Deactivate() {
	t.Status = StatusInactive
	t.UpdatedAt = time.Now()
}
