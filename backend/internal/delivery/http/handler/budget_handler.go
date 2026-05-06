package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/domain/finance"
	financeUC "github.com/protone/erp/internal/usecase/finance"
	"github.com/protone/erp/pkg/response"
)

type BudgetHandler struct {
	uc financeUC.UseCase
}

func NewBudgetHandler(uc financeUC.UseCase) *BudgetHandler {
	return &BudgetHandler{uc: uc}
}

func (h *BudgetHandler) SetBudget(w http.ResponseWriter, r *http.Request) {
	var body struct {
		DepartmentID string  `json:"department_id"`
		Month        int     `json:"month"`
		Year         int     `json:"year"`
		Amount       float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	deptID, _ := uuid.Parse(body.DepartmentID)

	budget := finance.DepartmentBudget{
		ID:              uuid.New(),
		TenantID:        tenantID,
		DepartmentID:    deptID,
		Month:           body.Month,
		Year:            body.Year,
		AllocatedAmount: body.Amount,
		SpentAmount:     0,
	}

	if err := h.uc.SetBudget(r.Context(), &budget); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "budget updated successfully", nil)
}
