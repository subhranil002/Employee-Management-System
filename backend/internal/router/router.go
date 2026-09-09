package router

import (
	"net/http"

	"github.com/subhranil002/GO-Cognito/internal/auth"
	"github.com/subhranil002/GO-Cognito/internal/employee"
	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

func Setup(
	authHandler *auth.Handler,
	tokenVerifier *middleware.TokenVerifier,
	employeeHandler *employee.Handler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, true, "ok", nil)
	})

	// Public routes
	mux.HandleFunc("GET /auth/status", authHandler.GetUserStatus)
	mux.HandleFunc("POST /auth/signup", authHandler.Signup)
	mux.HandleFunc("POST /auth/confirm", authHandler.ConfirmSignup)
	mux.HandleFunc("POST /auth/signin", authHandler.Signin)

	// Protected auth routes
	mux.HandleFunc("GET /profile", tokenVerifier.RequireAuth(authHandler.Profile))
	mux.HandleFunc("GET /auth/logout", tokenVerifier.RequireAuth(authHandler.Logout))

	// Protected employee proxy routes
	mux.HandleFunc("GET /employees", tokenVerifier.RequireAuth(employeeHandler.List))
	mux.HandleFunc("POST /employees", tokenVerifier.RequireAuth(employeeHandler.Create))
	mux.HandleFunc("GET /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Get))
	mux.HandleFunc("PATCH /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Update))
	mux.HandleFunc("DELETE /employees/{id}", tokenVerifier.RequireAuth(employeeHandler.Delete))

	return mux
}
