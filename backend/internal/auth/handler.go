package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/subhranil002/GO-Cognito/internal/middleware"
	"github.com/subhranil002/GO-Cognito/pkg/cookies"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest

	// Parse and validate request body
	if err := readJSON(w, r, &req); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Register user in Cognito
	result, err := h.service.Signup(ctx, req)
	if err != nil {
		handleError(w, err)
		return
	}

	response.JSON(
		w,
		http.StatusCreated,
		true,
		"signup successful; check your email for the confirmation code",
		result,
	)
}

func (h *Handler) ConfirmSignup(w http.ResponseWriter, r *http.Request) {
	var req ConfirmSignUpRequest

	// Parse and validate request body
	if err := readJSON(w, r, &req); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Confirm user account with code sent via email
	if err := h.service.ConfirmSignup(ctx, req); err != nil {
		handleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, true, "account confirmed", nil)
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest

	// Parse and validate request body
	if err := readJSON(w, r, &req); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Authenticate credentials with Cognito
	result, err := h.service.Signin(ctx, req)
	if err != nil {
		handleError(w, err)
		return
	}

	// Store tokens in secure HttpOnly cookies
	cookies.SetAuthCookies(w, result.AccessToken, result.RefreshToken, int(result.ExpiresIn))
	response.JSON(w, http.StatusOK, true, "signin successful", nil)
}

func (h *Handler) Profile(
	w http.ResponseWriter,
	r *http.Request,
	accessToken string,
	claims *middleware.AccessTokenClaims,
) {
	// Fetch complete user profile from Cognito
	profile, err := h.service.GetProfile(r.Context(), accessToken)
	if err != nil {
		handleError(w, err)
		return
	}

	// Send verified user profile back to client
	response.JSON(w, http.StatusOK, true, "profile fetched successfully", profile)
}

func (h *Handler) Logout(
	w http.ResponseWriter,
	r *http.Request,
	accessToken string,
	_ *middleware.AccessTokenClaims,
) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Invalidate user sessions globally in Cognito
	if err := h.service.Logout(ctx, accessToken); err != nil {
		handleError(w, err)
		return
	}

	// Remove authentication cookies from client
	cookies.ClearAuthCookies(w)
	response.JSON(w, http.StatusOK, true, "logged out", nil)
}

func (h *Handler) GetUserStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	email := strings.TrimSpace(
		r.URL.Query().Get("email"),
	)

	if email == "" {
		response.Error(
			w,
			http.StatusBadRequest,
			"email is required",
		)
		return
	}

	ctx, cancel := context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
	defer cancel()

	status, err := h.service.GetUserStatus(
		ctx,
		email,
	)

	if err != nil {
		slog.Error(
			"failed to get user status",
			"error", err,
		)

		response.Error(
			w,
			http.StatusInternalServerError,
			"unable to check user status",
		)
		return
	}

	response.JSON(
		w,
		http.StatusOK,
		true,
		"user status fetched",
		UserStatusResponse{
			Status: status,
		},
	)
}

func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Enforce payload size limit and reject unknown fields
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		slog.Error("failed to decode json request", "error", err, "path", r.URL.Path)
		response.Error(
			w,
			http.StatusBadRequest,
			fmt.Sprintf("invalid JSON: %v", err),
		)
		return err
	}

	return nil
}

func handleError(w http.ResponseWriter, err error) {
	slog.Error("request failed", "error", err)

	// Map internal sentinel errors to appropriate HTTP status responses
	switch {
	case errors.Is(err, ErrInvalid):
		response.Error(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusUnauthorized, err.Error())

	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
