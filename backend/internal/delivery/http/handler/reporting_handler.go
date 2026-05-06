package handler

import (
	"net/http"
	"strconv"

	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/reporting"
	"github.com/protone/erp/pkg/response"
)

type ReportingHandler struct {
	usecase *reporting.ReportingUseCase
}

func NewReportingHandler(uc *reporting.ReportingUseCase) *ReportingHandler {
	return &ReportingHandler{usecase: uc}
}

func (h *ReportingHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())

	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	if month == 0 || year == 0 {
		response.BadRequest(w, "Month and year are required")
		return
	}

	report, err := h.usecase.GetDashboardData(r.Context(), tenantID, month, year)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Dashboard stats fetched", report)
}
