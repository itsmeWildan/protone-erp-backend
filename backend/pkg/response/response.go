package response

import (
	"encoding/json"
	"net/http"
)

// ─── Standard Response Structures ───────────────────────────────────────────

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type ErrorDetail struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Errors  []ErrorDetail `json:"errors,omitempty"`
}

// ─── Helper Functions ────────────────────────────────────────────────────────

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Created(w http.ResponseWriter, message string, data any) {
	JSON(w, http.StatusCreated, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func SuccessPaginated(w http.ResponseWriter, message string, data any, page, perPage int, total int64) {
	JSON(w, http.StatusOK, Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:    page,
			PerPage: perPage,
			Total:   total,
		},
	})
}

func BadRequest(w http.ResponseWriter, message string, errs ...ErrorDetail) {
	JSON(w, http.StatusBadRequest, ErrorResponse{
		Status:  "error",
		Message: message,
		Errors:  errs,
	})
}

func Unauthorized(w http.ResponseWriter, message string) {
	JSON(w, http.StatusUnauthorized, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func Forbidden(w http.ResponseWriter, message string) {
	JSON(w, http.StatusForbidden, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func NotFound(w http.ResponseWriter, message string) {
	JSON(w, http.StatusNotFound, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func InternalServerError(w http.ResponseWriter, message string) {
	JSON(w, http.StatusInternalServerError, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}

func Conflict(w http.ResponseWriter, message string) {
	JSON(w, http.StatusConflict, ErrorResponse{
		Status:  "error",
		Message: message,
	})
}
