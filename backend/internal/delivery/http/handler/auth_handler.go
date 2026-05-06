package handler

import (
	"encoding/json"
	"net/http"

	"github.com/protone/erp/internal/usecase/auth"
	"github.com/protone/erp/pkg/response"
)

type AuthHandler struct {
	registerTenantUC *auth.RegisterTenantUseCase
	loginUC          *auth.LoginUseCase
}

func NewAuthHandler(registerTenantUC *auth.RegisterTenantUseCase, loginUC *auth.LoginUseCase) *AuthHandler {
	return &AuthHandler{
		registerTenantUC: registerTenantUC,
		loginUC:          loginUC,
	}
}

type registerTenantRequest struct {
	CompanyName string `json:"company_name"`
	Slug        string `json:"slug"`
	AdminEmail  string `json:"admin_email"`
	Password    string `json:"password"`
}

func (h *AuthHandler) RegisterTenant(w http.ResponseWriter, r *http.Request) {
	var req registerTenantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	input := auth.RegisterTenantInput{
		CompanyName: req.CompanyName,
		Slug:        req.Slug,
		AdminEmail:  req.AdminEmail,
		Password:    req.Password,
	}

	out, err := h.registerTenantUC.Execute(r.Context(), input)
	if err != nil {
		// Logika handling error domain bisa ditambahkan di sini
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, "Tenant and Admin created successfully", out)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	out, err := h.loginUC.Execute(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		response.Unauthorized(w, "Invalid credentials")
		return
	}

	response.Success(w, "Login successful", out)
}
