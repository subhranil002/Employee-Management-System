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

// List fetches all employees for the authenticated user.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.UserIDFromContext(r.Context())

	employees, err := h.client.List(r.Context(), userID)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "employees retrieved successfully", employees)
}

// Get fetches a single employee by ID for the authenticated user.
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

// Create creates a new employee for the authenticated user.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEmployeeRequest

	if err := readJSON(w, r, &req); err != nil {
		return
	}

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

// Update applies a partial update to an existing employee for the authenticated user.
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

// Delete removes an employee by ID for the authenticated user.
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

// readJSON parses and validates the JSON request body.
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

// handleError maps backend errors to appropriate HTTP responses.
func handleError(w http.ResponseWriter, err error) {
	slog.Error("employee request failed", "error", err)

	var be *backendError
	if errors.As(err, &be) {
		response.Error(w, be.StatusCode, be.Message)
		return
	}

	response.Error(w, http.StatusInternalServerError, "internal server error")
}
