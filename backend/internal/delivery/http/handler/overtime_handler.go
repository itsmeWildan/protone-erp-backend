package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/overtime"
	"github.com/protone/erp/pkg/response"
)

type OvertimeHandler struct {
	usecase *overtime.OvertimeUseCase
}

func NewOvertimeHandler(uc *overtime.OvertimeUseCase) *OvertimeHandler {
	return &OvertimeHandler{usecase: uc}
}

func (h *OvertimeHandler) Request(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Date          string  `json:"date"`
		StartTime     string  `json:"start_time"`
		EndTime       string  `json:"end_time"`
		DurationHours float64 `json:"duration_hours"`
		Reason        string  `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	date, _ := time.Parse("2006-01-02", body.Date)

	input := overtime.SubmitOvertimeInput{
		TenantID:      tenantID,
		UserID:        userID,
		Date:          date,
		StartTime:     body.StartTime,
		EndTime:       body.EndTime,
		DurationHours: body.DurationHours,
		Reason:        body.Reason,
	}

	if err := h.usecase.SubmitRequest(r.Context(), input); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, "Overtime request submitted", nil)
}

func (h *OvertimeHandler) GetMyRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))

	list, err := h.usecase.GetMyOvertime(r.Context(), tenantID, userID, month, year)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Overtime requests fetched", list)
}

func (h *OvertimeHandler) Approve(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RequestID string `json:"request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())
	reqID, _ := uuid.Parse(body.RequestID)

	if err := h.usecase.ApproveRequest(r.Context(), tenantID, userID, reqID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Overtime request approved", nil)
}
