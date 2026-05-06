package jwt

import (
	"fmt"
	"time"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/protone/erp/config"
)

// Claims adalah payload JWT kita.
// Berisi user_id, tenant_id, dan role — semua yang dibutuhkan middleware.
type Claims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	gojwt.RegisteredClaims
}

type Manager struct {
	cfg config.JWTConfig
}

func NewManager(cfg config.JWTConfig) *Manager {
	return &Manager{cfg: cfg}
}

// GenerateAccessToken membuat JWT access token
func (m *Manager) GenerateAccessToken(userID, tenantID uuid.UUID, role string) (string, error) {
	claims := Claims{
		UserID:   userID.String(),
		TenantID: tenantID.String(),
		Role:     role,
		RegisteredClaims: gojwt.RegisteredClaims{
			ExpiresAt: gojwt.NewNumericDate(time.Now().Add(time.Duration(m.cfg.AccessTTLMinutes) * time.Minute)),
			IssuedAt:  gojwt.NewNumericDate(time.Now()),
			Issuer:    "protone-erp",
		},
	}

	token := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.cfg.Secret))
}

// GenerateRefreshToken membuat token acak untuk refresh
func (m *Manager) GenerateRefreshToken() (string, error) {
	id := uuid.New()
	return id.String(), nil
}

// VerifyToken memvalidasi JWT dan extract claims
func (m *Manager) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := gojwt.ParseWithClaims(tokenStr, &Claims{}, func(token *gojwt.Token) (any, error) {
		if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
