package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/protone/erp/internal/domain/user"
	"github.com/protone/erp/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginUseCase struct {
	userRepo   user.Repository
	jwtManager *jwt.Manager
}

func NewLogin(ur user.Repository, jm *jwt.Manager) *LoginUseCase {
	return &LoginUseCase{
		userRepo:   ur,
		jwtManager: jm,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// 1. Get user by email
	u, err := uc.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("invalid email or password")
	}

	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(input.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// 3. Generate tokens
	accessToken, err := uc.jwtManager.GenerateAccessToken(u.ID, u.TenantID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.jwtManager.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
