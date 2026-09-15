package router

import (
	"net/http"

	"github.com/subhranil002/GO-Cognito/internal/employee"
	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

// Setup registers all application routes and returns the configured ServeMux.
// All employee routes are protected by JWT authentication via RequireAuth.
func Setup(
	tokenVerifier *middleware.TokenVerifier,
	employeeHandler *employee.Handler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Health check (public)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, true, "ok", nil)
	})

	// Employee CRUD routes (JWT-protected)
	mux.HandleFunc("GET /employees", tokenVerifier.RequireAuth(employeeHandler.List))
	mux.HandleFunc("POST /employees", tokenVerifier.RequireAuth(employeeHandler.Create))
	mux.HandleFunc("GET /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Get))
	mux.HandleFunc("PATCH /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Update))
	mux.HandleFunc("DELETE /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Delete))

	return mux
}
