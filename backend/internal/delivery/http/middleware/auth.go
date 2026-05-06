package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	pkgjwt "github.com/protone/erp/pkg/jwt"
	"github.com/protone/erp/pkg/response"
)

// ─── Context Keys ────────────────────────────────────────────────────────────

type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyTenantID contextKey = "tenant_id"
	ContextKeyRole     contextKey = "role"
)

// ─── Auth Middleware ─────────────────────────────────────────────────────────

type AuthMiddleware struct {
	jwtManager *pkgjwt.Manager
}

func NewAuthMiddleware(jwtManager *pkgjwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwtManager: jwtManager}
}

// Authenticate memvalidasi JWT dan inject user_id + tenant_id + role ke context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.URL.Query().Get("token")
		authHeader := r.Header.Get("Authorization")
		
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			response.Unauthorized(w, "Authorization token is required")
			return
		}

		claims, err := m.jwtManager.VerifyToken(tokenString)
		if err != nil {
			response.Unauthorized(w, "Invalid or expired token")
			return
		}

		// Inject ke context
		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyTenantID, claims.TenantID)
		ctx = context.WithValue(ctx, ContextKeyRole, claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ─── RBAC Middleware ─────────────────────────────────────────────────────────

// RequireRole memastikan user memiliki role yang diperlukan
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok {
				response.Forbidden(w, "Access denied")
				return
			}

			for _, allowed := range roles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Forbidden(w, "You don't have permission to access this resource")
		})
	}
}

// ─── Context Helpers ─────────────────────────────────────────────────────────

func GetTenantID(ctx context.Context) (uuid.UUID, bool) {
	s, ok := ctx.Value(ContextKeyTenantID).(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	s, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok {
		return uuid.Nil, false
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}
