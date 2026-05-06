package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/reimbursement"
	"github.com/protone/erp/pkg/response"
)

type ReimbursementHandler struct {
	uc reimbursement.UseCase
}

func NewReimbursementHandler(uc reimbursement.UseCase) *ReimbursementHandler {
	return &ReimbursementHandler{uc: uc}
}

func (h *ReimbursementHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Category    string  `json:"category"`
		Amount      float64 `json:"amount"`
		Date        string  `json:"date"`
		Description string  `json:"description"`
		ReceiptURL  string  `json:"receipt_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	date, err := time.Parse("2006-01-02", body.Date)
	if err != nil {
		date = time.Now()
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	input := reimbursement.SubmitClaimInput{
		TenantID:    tenantID,
		UserID:      userID,
		Category:    body.Category,
		Amount:      body.Amount,
		Date:        date,
		Description: body.Description,
	}
	if body.ReceiptURL != "" {
		input.ReceiptURL = &body.ReceiptURL
	}

	if err := h.uc.SubmitClaim(r.Context(), input); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, "Reimbursement claim submitted successfully", nil)
}

func (h *ReimbursementHandler) GetMyClaims(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	claims, err := h.uc.GetMyClaims(r.Context(), tenantID, userID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "My claims fetched", claims)
}

func (h *ReimbursementHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ClaimID string `json:"claim_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	claimID, err := uuid.Parse(body.ClaimID)
	if err != nil {
		response.BadRequest(w, "invalid claim id")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	if err := h.uc.ApproveClaim(r.Context(), tenantID, userID, claimID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Reimbursement claim approved successfully", nil)
}
