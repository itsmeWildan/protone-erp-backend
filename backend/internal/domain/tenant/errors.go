package tenant

import "errors"

var (
	ErrNameRequired   = errors.New("tenant name is required")
	ErrSlugRequired   = errors.New("tenant slug is required")
	ErrInvalidStatus  = errors.New("invalid tenant status")
	ErrTenantNotFound = errors.New("tenant not found")
)
