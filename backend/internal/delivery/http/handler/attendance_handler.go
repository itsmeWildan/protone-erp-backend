package handler

import (
	"encoding/json"
	"net/http"

	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/attendance"
	"github.com/protone/erp/pkg/response"
)

type AttendanceHandler struct {
	uc attendance.UseCase
}

func NewAttendanceHandler(uc attendance.UseCase) *AttendanceHandler {
	return &AttendanceHandler{uc: uc}
}

func (h *AttendanceHandler) ClockIn(w http.ResponseWriter, r *http.Request) {
	var body struct {
		LocationIn string `json:"location_in"`
		Notes      string `json:"notes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	input := attendance.ClockInInput{
		TenantID:   tenantID,
		UserID:     userID,
		LocationIn: body.LocationIn,
		Notes:      body.Notes,
	}

	if err := h.uc.ClockIn(r.Context(), input); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Clock-in successful", nil)
}

func (h *AttendanceHandler) ClockOut(w http.ResponseWriter, r *http.Request) {
	var body struct {
		LocationOut string `json:"location_out"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	input := attendance.ClockOutInput{
		TenantID:    tenantID,
		UserID:      userID,
		LocationOut: body.LocationOut,
	}

	if err := h.uc.ClockOut(r.Context(), input); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Clock-out successful", nil)
}
