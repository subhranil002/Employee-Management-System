package router

import (
	"net/http"

	"github.com/subhranil002/GO-Cognito/internal/employee"
	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/internal/user"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

// Setup registers application routes and returns the configured ServeMux
func Setup(
	verifier *middleware.Verifier,
	employeeHandler *employee.Handler,
	userHandler *user.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()

	// Public health check route
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, true, "ok", nil)
	})

	// Protected user profile route
	mux.HandleFunc("GET /profile", verifier.RequireAuth(userHandler.Profile))

	// Protected employee CRUD routes
	mux.HandleFunc("GET /employees", verifier.RequireAuth(employeeHandler.List))
	mux.HandleFunc("POST /employees", verifier.RequireAuth(employeeHandler.Create))
	mux.HandleFunc("GET /employees/{id}", verifier.RequireAuth(employeeHandler.Get))
	mux.HandleFunc("PATCH /employees/{id}", verifier.RequireAuth(employeeHandler.Update))
	mux.HandleFunc("DELETE /employees/{id}", verifier.RequireAuth(employeeHandler.Delete))

	return mux
}
