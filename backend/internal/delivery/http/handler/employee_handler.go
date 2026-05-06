package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/domain/employee"
	empcommand "github.com/protone/erp/internal/usecase/employee/command"
	empquery "github.com/protone/erp/internal/usecase/employee/query"
	"github.com/protone/erp/pkg/response"
)

type EmployeeHandler struct {
	createUC *empcommand.CreateEmployeeUseCase
	updateUC *empcommand.UpdateEmployeeUseCase
	deleteUC *empcommand.DeleteEmployeeUseCase
	listUC   *empquery.ListEmployeesUseCase
	getUC    *empquery.GetEmployeeUseCase
}

func NewEmployeeHandler(
	createUC *empcommand.CreateEmployeeUseCase,
	updateUC *empcommand.UpdateEmployeeUseCase,
	deleteUC *empcommand.DeleteEmployeeUseCase,
	listUC *empquery.ListEmployeesUseCase,
	getUC *empquery.GetEmployeeUseCase,
) *EmployeeHandler {
	return &EmployeeHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		listUC:   listUC,
		getUC:    getUC,
	}
}

// ─── List Employees ─────────────────────────────────────────────────────────
// GET /api/v1/employees

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.GetTenantID(r.Context())
	if !ok {
		response.Unauthorized(w, "Invalid tenant")
		return
	}

	q := r.URL.Query()
	input := empquery.ListEmployeesInput{
		TenantID: tenantID,
		Search:   q.Get("search"),
		Page:     parseIntQuery(q.Get("page"), 1),
		PerPage:  parseIntQuery(q.Get("per_page"), 20),
	}

	if deptStr := q.Get("department_id"); deptStr != "" {
		if id, err := uuid.Parse(deptStr); err == nil {
			input.DepartmentID = &id
		}
	}
	if status := q.Get("status"); status != "" {
		input.Status = &status
	}

	out, err := h.listUC.Execute(r.Context(), input)
	if err != nil {
		response.InternalServerError(w, "Failed to fetch employees")
		return
	}

	response.SuccessPaginated(w, "Employees fetched successfully",
		out.Items, out.Page, out.PerPage, out.Total)
}

// ─── Get Employee ────────────────────────────────────────────────────────────
// GET /api/v1/employees/{id}

func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.GetTenantID(r.Context())
	if !ok {
		response.Unauthorized(w, "Invalid tenant")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid employee ID")
		return
	}

	out, err := h.getUC.Execute(r.Context(), tenantID, id)
	if err != nil {
		if errors.Is(err, employee.ErrNotFound) {
			response.NotFound(w, "Employee not found")
			return
		}
		response.InternalServerError(w, "Failed to fetch employee")
		return
	}

	response.Success(w, "Employee fetched successfully", out)
}

// ─── Create Employee ─────────────────────────────────────────────────────────
// POST /api/v1/employees

type createEmployeeRequest struct {
	NIK            string  `json:"nik"`
	FullName       string  `json:"full_name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Gender         string  `json:"gender"`
	BirthDate      string  `json:"birth_date"`           // "YYYY-MM-DD"
	DepartmentID   string  `json:"department_id"`
	PositionID     string  `json:"position_id"`
	ManagerID      string  `json:"manager_id"`
	EmploymentType string  `json:"employment_type"`
	JoinDate       string  `json:"join_date"`             // "YYYY-MM-DD"
	BasicSalary    float64 `json:"basic_salary"`
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.GetTenantID(r.Context())
	if !ok {
		response.Unauthorized(w, "Invalid tenant")
		return
	}

	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	input := empcommand.CreateEmployeeInput{
		TenantID:       tenantID,
		NIK:            req.NIK,
		FullName:       req.FullName,
		Email:          req.Email,
		Phone:          req.Phone,
		Gender:         req.Gender,
		EmploymentType: req.EmploymentType,
		BasicSalary:    req.BasicSalary,
	}

	// Parse dates
	if req.JoinDate != "" {
		t, err := time.Parse("2006-01-02", req.JoinDate)
		if err != nil {
			response.BadRequest(w, "Invalid join_date format (use YYYY-MM-DD)")
			return
		}
		input.JoinDate = t
	}
	if req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			response.BadRequest(w, "Invalid birth_date format (use YYYY-MM-DD)")
			return
		}
		input.BirthDate = &t
	}

	// Parse UUIDs
	if id, err := uuid.Parse(req.DepartmentID); err == nil {
		input.DepartmentID = id
	}
	if id, err := uuid.Parse(req.PositionID); err == nil {
		input.PositionID = id
	}
	if req.ManagerID != "" {
		if id, err := uuid.Parse(req.ManagerID); err == nil {
			input.ManagerID = &id
		}
	}

	out, err := h.createUC.Execute(r.Context(), input)
	if err != nil {
		if errors.Is(err, employee.ErrNIKAlreadyExists) {
			response.Conflict(w, "NIK already exists for this company")
			return
		}
		if errors.Is(err, employee.ErrNIKRequired) ||
			errors.Is(err, employee.ErrFullNameRequired) ||
			errors.Is(err, employee.ErrJoinDateRequired) {
			response.BadRequest(w, err.Error())
			return
		}
		response.InternalServerError(w, "Failed to create employee")
		return
	}

	response.Created(w, "Employee created successfully", map[string]string{"id": out.ID})
}

// ─── Update Employee ─────────────────────────────────────────────────────────
// PUT /api/v1/employees/{id}

type updateEmployeeRequest struct {
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	DepartmentID    string  `json:"department_id"`
	PositionID      string  `json:"position_id"`
	ManagerID       string  `json:"manager_id"`
	BasicSalary     float64 `json:"basic_salary"`
	BankName        string  `json:"bank_name"`
	BankAccountNo   string  `json:"bank_account_no"`
	BankAccountName string  `json:"bank_account_name"`
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.GetTenantID(r.Context())
	if !ok {
		response.Unauthorized(w, "Invalid tenant")
		return
	}

	empID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid employee ID")
		return
	}

	var req updateEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	input := empcommand.UpdateEmployeeInput{
		TenantID:        tenantID,
		ID:              empID,
		FullName:        req.FullName,
		Email:           req.Email,
		Phone:           req.Phone,
		BasicSalary:     req.BasicSalary,
		BankName:        req.BankName,
		BankAccountNo:   req.BankAccountNo,
		BankAccountName: req.BankAccountName,
	}

	if id, err := uuid.Parse(req.DepartmentID); err == nil {
		input.DepartmentID = id
	}
	if id, err := uuid.Parse(req.PositionID); err == nil {
		input.PositionID = id
	}
	if req.ManagerID != "" {
		if id, err := uuid.Parse(req.ManagerID); err == nil {
			input.ManagerID = &id
		}
	}

	if err := h.updateUC.Execute(r.Context(), input); err != nil {
		if errors.Is(err, employee.ErrNotFound) {
			response.NotFound(w, "Employee not found")
			return
		}
		response.InternalServerError(w, "Failed to update employee")
		return
	}

	response.Success(w, "Employee updated successfully", nil)
}

// ─── Delete Employee ─────────────────────────────────────────────────────────
// DELETE /api/v1/employees/{id}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := middleware.GetTenantID(r.Context())
	if !ok {
		response.Unauthorized(w, "Invalid tenant")
		return
	}

	empID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "Invalid employee ID")
		return
	}

	if err := h.deleteUC.Execute(r.Context(), tenantID, empID); err != nil {
		if errors.Is(err, employee.ErrNotFound) {
			response.NotFound(w, "Employee not found")
			return
		}
		response.InternalServerError(w, "Failed to delete employee")
		return
	}

	response.Success(w, "Employee deleted successfully", nil)
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func parseIntQuery(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
