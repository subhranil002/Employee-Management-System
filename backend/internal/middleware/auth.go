package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/subhranil002/GO-Cognito/pkg/cookies"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

type AuthenticatedHandlerFunc func(
	w http.ResponseWriter,
	r *http.Request,
	accessToken string,
	claims *AccessTokenClaims,
)

func (v *TokenVerifier) RequireAuth(next AuthenticatedHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract access token from cookie with fallback to Authorization header
		var accessToken string
		if cookie, err := r.Cookie("access_token"); err == nil {
			accessToken = cookie.Value
		}
		if accessToken == "" {
			header := r.Header.Get("Authorization")
			parts := strings.Fields(header)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				accessToken = parts[1]
			}
		}

		if accessToken == "" {
			slog.Error("missing access token", "path", r.URL.Path)
			response.Error(
				w,
				http.StatusUnauthorized,
				"missing access token",
			)
			return
		}

		// Verify access token signature and claims
		claims, err := v.Verify(accessToken)
		if err != nil {
			// Auto-refresh access token if expired and refresh token is provided
			if errors.Is(err, jwt.ErrTokenExpired) && v.refresher != nil {
				var refreshToken string
				if cookie, err := r.Cookie("refresh_token"); err == nil {
					refreshToken = cookie.Value
				} else {
					refreshToken = r.Header.Get("X-Refresh-Token")
				}

				if refreshToken != "" {
					// Parse expired claims to retrieve username for Cognito SECRET_HASH
					expiredClaims, parseErr := parseExpiredClaims(accessToken)

					if parseErr == nil {
						// Exchange refresh token for a new access token
						newTokens, refreshErr := v.refresher.Refresh(
							r.Context(),
							refreshToken,
							expiredClaims.Username,
						)

						if refreshErr == nil {
							// Verify newly refreshed token
							newClaims, verifyErr := v.Verify(newTokens.AccessToken)

							if verifyErr == nil {
								// Update browser cookies with the new access token
								cookies.SetAuthCookies(
									w,
									newTokens.AccessToken,
									refreshToken,
									int(newTokens.ExpiresIn),
								)

								// Proceed to downstream handler with new token and claims
								next(w, r, newTokens.AccessToken, newClaims)
								return
							} else {
								slog.Error("failed to verify refreshed token", "error", verifyErr)
							}
						} else {
							slog.Error("failed to refresh token", "error", refreshErr)
						}
					} else {
						slog.Error("failed to parse expired claims", "error", parseErr)
					}
				}
			}

			slog.Error("invalid or expired access token", "error", err, "path", r.URL.Path)
			response.Error(
				w,
				http.StatusUnauthorized,
				"invalid or expired access token",
			)
			return
		}

		// Pass verified token and claims to downstream handler
		next(w, r, accessToken, claims)
	}
}
