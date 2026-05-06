package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

type User struct {
	ID           uuid.UUID
	TenantID     uuid.UUID
	Name         string
	Email        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

func NewUser(tenantID uuid.UUID, name, email, password string, role Role) (*User, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           uuid.New(),
		TenantID:     tenantID,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}, nil
}
