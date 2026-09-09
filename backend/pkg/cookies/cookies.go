package cookies

import (
	"net/http"
	"time"
)

func SetAuthCookies(w http.ResponseWriter, accessToken, refreshToken string, expiresIn int) {
	// Set HttpOnly access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   expiresIn,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// Set HttpOnly refresh token cookie when present
	if refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 30, // 30 days
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}

func ClearAuthCookies(w http.ResponseWriter) {
	cookiesToClear := []string{"access_token", "refresh_token"}

	// Expire authentication cookies immediately
	for _, name := range cookiesToClear {
		http.SetCookie(w, &http.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	}
}
