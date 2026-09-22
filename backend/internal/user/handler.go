package user

import (
	"net/http"

	"github.com/subhranil002/GO-Cognito/pkg/response"
)

// Handler handles user-related HTTP requests
type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// Profile returns the authenticated user profile from context
func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	u, ok := FromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	response.JSON(w, http.StatusOK, true, "profile retrieved successfully", u)
}
