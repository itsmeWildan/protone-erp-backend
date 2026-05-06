package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/protone/erp/internal/delivery/http/middleware"
	"github.com/protone/erp/internal/usecase/leave"
	"github.com/protone/erp/pkg/response"
)

type LeaveHandler struct {
	usecase *leave.LeaveUseCase
}

func NewLeaveHandler(uc *leave.LeaveUseCase) *LeaveHandler {
	return &LeaveHandler{usecase: uc}
}

func (h *LeaveHandler) GetTypes(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	types, err := h.usecase.GetLeaveTypes(r.Context(), tenantID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}
	response.Success(w, "Leave types fetched", types)
}

func (h *LeaveHandler) RequestLeave(w http.ResponseWriter, r *http.Request) {
	var body struct {
		LeaveTypeID string `json:"leave_type_id"`
		StartDate   string `json:"start_date"`
		EndDate     string `json:"end_date"`
		Reason      string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	ltID, _ := uuid.Parse(body.LeaveTypeID)
	start, _ := time.Parse("2006-01-02", body.StartDate)
	end, _ := time.Parse("2006-01-02", body.EndDate)

	input := leave.RequestLeaveInput{
		TenantID:    tenantID,
		UserID:      userID,
		LeaveTypeID: ltID,
		StartDate:   start,
		EndDate:     end,
		Reason:      body.Reason,
	}

	leaveID, err := h.usecase.RequestLeave(r.Context(), input)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Created(w, "Leave request submitted successfully", map[string]string{
		"leave_id": leaveID.String(),
	})
}

func (h *LeaveHandler) GetMyRequests(w http.ResponseWriter, r *http.Request) {
	tenantID, _ := middleware.GetTenantID(r.Context())
	userID, _ := middleware.GetUserID(r.Context())

	leaves, err := h.usecase.GetMyLeaves(r.Context(), tenantID, userID)
	if err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Leave requests fetched successfully", leaves)
}

func (h *LeaveHandler) Approve(w http.ResponseWriter, r *http.Request) {
	// Ambil ID cuti dari query param atau body (kali ini kita pakai body sederhana)
	var body struct {
		LeaveID string `json:"leave_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	managerUserID, _ := middleware.GetUserID(r.Context())
	leaveID, _ := uuid.Parse(body.LeaveID)

	if err := h.usecase.ApproveLeave(r.Context(), tenantID, managerUserID, leaveID); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Leave request approved successfully", nil)
}
func (h *LeaveHandler) GiveBalance(w http.ResponseWriter, r *http.Request) {
	var body struct {
		EmployeeID  string `json:"employee_id"`
		LeaveTypeID string `json:"leave_type_id"`
		Year        int    `json:"year"`
		TotalDays   int    `json:"total_days"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	tenantID, _ := middleware.GetTenantID(r.Context())
	empID, _ := uuid.Parse(body.EmployeeID)
	ltID, _ := uuid.Parse(body.LeaveTypeID)

	if err := h.usecase.GiveBalance(r.Context(), tenantID, empID, ltID, body.Year, body.TotalDays); err != nil {
		response.InternalServerError(w, err.Error())
		return
	}

	response.Success(w, "Leave balance given successfully", nil)
}
