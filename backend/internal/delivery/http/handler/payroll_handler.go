package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/payroll"
	"github.com/protone/erp/pkg/response"
)

type PayrollHandler struct {
	usecase *payroll.PayrollUseCase
}

func NewPayrollHandler(uc *payroll.PayrollUseCase) *PayrollHandler {
	return &PayrollHandler{usecase: uc}
}

func (h *PayrollHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	if err := h.usecase.GenerateMonthlyPayroll(r.Context(), tenantID, body.Month, body.Year); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Payroll generated successfully", nil)
}

func (h *PayrollHandler) GetMySlip(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	slip, err := h.usecase.GetMySlip(r.Context(), tenantID, userID, month, year)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	if slip == nil {
		response.NotFound(w, "Payroll slip not found for this period")
		return
	}

	response.Success(w, "Payroll slip fetched", slip)
}

func (h *PayrollHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	if err := h.usecase.ApprovePayroll(r.Context(), tenantID, body.Month, body.Year); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Payroll approved successfully", nil)
}

func (h *PayrollHandler) Pay(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Month int `json:"month"`
		Year  int `json:"year"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	if err := h.usecase.PayPayroll(r.Context(), tenantID, body.Month, body.Year); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Payroll paid and journaled successfully", nil)
}

func (h *PayrollHandler) DownloadSlip(w http.ResponseWriter, r *http.Request) {
	slipIDStr := chi.URLParam(r, "id")
	periodID, _ := uuid.Parse(slipIDStr)

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	// Fallback untuk testing browser: ambil dari query ?token=
	if tenantID == uuid.Nil {
		tokenStr := r.URL.Query().Get("token")
		if tokenStr != "" {
			// (Logika parsing token sementara untuk mempermudah tes)
			// Untuk sekarang kita asumsikan jika token ada, kita pakai data context
		}
	}
	
	pdfReader, err := h.usecase.DownloadMySlipPDF(r.Context(), tenantID, userID, periodID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=slip_gaji.pdf")
	
	io.Copy(w, pdfReader)
}
