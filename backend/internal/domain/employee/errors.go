package employee

import "errors"

var (
	ErrNIKRequired       = errors.New("NIK is required")
	ErrFullNameRequired  = errors.New("full name is required")
	ErrTenantIDRequired  = errors.New("tenant ID is required")
	ErrDepartmentRequired = errors.New("department is required")
	ErrPositionRequired  = errors.New("position is required")
	ErrJoinDateRequired  = errors.New("join date is required")
	ErrInvalidSalary     = errors.New("salary cannot be negative")
	ErrAlreadyTerminated = errors.New("employee is already terminated")
	ErrNotFound          = errors.New("employee not found")
	ErrNIKAlreadyExists  = errors.New("NIK already exists for this company")
)
