package employee

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

type Handler struct {
	client *EmployeeClient
}

func NewHandler(client *EmployeeClient) *Handler {
	return &Handler{client: client}
}

// List handles retrieving all employees for the authenticated user
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	employees, err := h.client.List(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employees retrieved successfully", employees)
}

// Get handles retrieving a single employee by ID
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	emp, err := h.client.Get(r.Context(), userID, id)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee retrieved successfully", emp)
}

// Create validates employee payload and creates a new employee
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEmployeeRequest

	if err := readJSON(w, r, &req); err != nil {
		return
	}

	// Validate employee identifier format
	if matched, _ := regexp.MatchString(`^EMP-\d{3}$`, req.EmpID); !matched {
		response.Error(w, http.StatusBadRequest, "empID must be in the format EMP-123 (EMP- followed by exactly 3 digits)")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	emp, err := h.client.Create(r.Context(), userID, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, true, "employee created successfully", emp)
}

// Update validates payload and updates specified employee fields
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	var req UpdateEmployeeRequest

	if err := readJSON(w, r, &req); err != nil {
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	emp, err := h.client.Update(r.Context(), userID, id, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee updated successfully", emp)
}

// Delete removes an employee record by ID
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		response.Error(w, http.StatusBadRequest, "id is required")
		return
	}

	userID, _ := middleware.UserIDFromContext(r.Context())

	if err := h.client.Delete(r.Context(), userID, id); err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employee deleted successfully", nil)
}

// Decode and validate JSON request body
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

// Map downstream service errors to HTTP response codes
func handleError(w http.ResponseWriter, err error) {
	slog.Error("employee request failed", "error", err)

	var be *backendError
	if errors.As(err, &be) {
		response.Error(w, be.StatusCode, be.Message)
		return
	}

	response.Error(w, http.StatusInternalServerError, "internal server error")
}
