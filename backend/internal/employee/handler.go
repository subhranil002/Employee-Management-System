package employee

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

type Handler struct {
	client *EmployeeClient
}

func NewHandler(client *EmployeeClient) *Handler {
	return &Handler{client: client}
}

// Fetch all employees
func (h *Handler) List(w http.ResponseWriter, r *http.Request, _ string, _ *middleware.AccessTokenClaims) {
	employees, err := h.client.List(r.Context())
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employees retrieved successfully", employees)
}

// Fetch single employee by ID
func (h *Handler) Get(w http.ResponseWriter, r *http.Request, _ string, _ *middleware.AccessTokenClaims) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	emp, err := h.client.Get(r.Context(), id)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee retrieved successfully", emp)
}

// Create new employee
func (h *Handler) Create(w http.ResponseWriter, r *http.Request, _ string, _ *middleware.AccessTokenClaims) {
	var req CreateEmployeeRequest

	if err := readJSON(w, r, &req); err != nil {
		return
	}

	emp, err := h.client.Create(r.Context(), req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, true, "employee created successfully", emp)
}

// Update existing employee
func (h *Handler) Update(w http.ResponseWriter, r *http.Request, _ string, _ *middleware.AccessTokenClaims) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateEmployeeRequest

	if err := readJSON(w, r, &req); err != nil {
		return
	}

	emp, err := h.client.Update(r.Context(), id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee updated successfully", emp)
}

// Delete employee by ID
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, _ string, _ *middleware.AccessTokenClaims) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.client.Delete(r.Context(), id); err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee deleted successfully", nil)
}

// Parse and validate request JSON body
func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		slog.Error("failed to decode employee request", "error", err, "path", r.URL.Path)
		response.Error(w, http.StatusBadRequest, fmt.Sprintf("invalid JSON: %v", err))
		return err
	}

	return nil
}

// Handle error response to client
func handleError(w http.ResponseWriter, err error) {
	slog.Error("employee request failed", "error", err)

	var be *backendError
	if errors.As(err, &be) {
		// Forward status code and message from CRUD backend
		response.Error(w, be.StatusCode, be.Message)
		return
	}

	response.Error(w, http.StatusInternalServerError, "internal server error")
}
