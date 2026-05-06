package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/protone/erp/internal/domain/tenant"
	"github.com/protone/erp/internal/domain/user"
	"github.com/protone/erp/internal/infrastructure/persistence/postgres"
)

type RegisterTenantInput struct {
	CompanyName string
	Slug        string
	AdminEmail  string
	Password    string
}

type RegisterTenantOutput struct {
	TenantID string
	AdminID  string
}

type RegisterTenantUseCase struct {
	tenantRepo tenant.Repository
	userRepo   user.Repository
	txManager  *postgres.TxManager
}

func NewRegisterTenant(tr tenant.Repository, ur user.Repository, tm *postgres.TxManager) *RegisterTenantUseCase {
	return &RegisterTenantUseCase{
		tenantRepo: tr,
		userRepo:   ur,
		txManager:  tm,
	}
}

func (uc *RegisterTenantUseCase) Execute(ctx context.Context, input RegisterTenantInput) (*RegisterTenantOutput, error) {
	// 1. Check if slug exists
	existingTenant, err := uc.tenantRepo.GetBySlug(ctx, input.Slug)
	if err != nil {
		return nil, err
	}
	if existingTenant != nil {
		return nil, errors.New("slug already taken")
	}

	// 2. Check if admin email exists
	existingUser, err := uc.userRepo.GetByEmail(ctx, input.AdminEmail)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// 3. Start Transaction
	var out RegisterTenantOutput
	err = uc.txManager.WithTx(ctx, func(txCtx context.Context) error {
		// Create Tenant entity
		t, err := tenant.NewTenant(tenant.CreateParams{
			Name:  input.CompanyName,
			Slug:  input.Slug,
			Email: input.AdminEmail,
		})
		if err != nil {
			return err
		}

		// Persist Tenant
		if err := uc.tenantRepo.Create(txCtx, t); err != nil {
			return err
		}

		// Create Admin User entity
		u, err := user.NewUser(t.ID, "Administrator", input.AdminEmail, input.Password, user.RoleAdmin)
		if err != nil {
			return err
		}

		// Persist User
		if err := uc.userRepo.Create(txCtx, u); err != nil {
			return err
		}

		out.TenantID = t.ID.String()
		out.AdminID = u.ID.String()
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("register tenant transaction: %w", err)
	}

	return &out, nil
}
