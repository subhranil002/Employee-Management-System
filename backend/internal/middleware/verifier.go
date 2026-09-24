package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/subhranil002/GO-Cognito/internal/user"
	"github.com/subhranil002/GO-Cognito/pkg/response"
)

type contextKey string

const (
	claimsContextKey contextKey = "cognito_claims"
	userIDContextKey contextKey = "user_id"
)

type AccessTokenClaims struct {
	TokenUse string `json:"token_use"`
	ClientID string `json:"client_id"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}

func ClaimsFromContext(ctx context.Context) (*AccessTokenClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*AccessTokenClaims)
	return claims, ok
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDContextKey).(string)
	return v, ok && v != ""
}

func UserFromContext(ctx context.Context) (*user.User, bool) {
	return user.FromContext(ctx)
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Convert base64url modulus and exponent components into an RSA public key
func parseRSAPublicKey(nB64, eB64 string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nB64)
	if err != nil {
		return nil, fmt.Errorf("decode modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eB64)
	if err != nil {
		return nil, fmt.Errorf("decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}

type Verifier struct {
	jwksURL    string
	issuer     string
	clientID   string
	httpClient *http.Client
	userClient *user.UserClient
}

func NewVerifier(issuer, clientID string, userClient *user.UserClient) *Verifier {
	return &Verifier{
		jwksURL:    issuer + "/.well-known/jwks.json",
		issuer:     issuer,
		clientID:   clientID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		userClient: userClient,
	}
}

// Retrieve and cache Cognito public verification keys
func (v *Verifier) fetchJWKS(ctx context.Context) (map[string]*rsa.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURL, nil)
	if err != nil {
		return nil, fmt.Errorf("build JWKS request: %w", err)
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JWKS endpoint returned HTTP %d", resp.StatusCode)
	}

	var raw jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("decode JWKS: %w", err)
	}

	// Parse RSA signing keys by key ID
	keys := make(map[string]*rsa.PublicKey, len(raw.Keys))
	for _, k := range raw.Keys {
		if k.Kty != "RSA" || k.Use != "sig" {
			continue
		}
		pub, err := parseRSAPublicKey(k.N, k.E)
		if err != nil {
			return nil, fmt.Errorf("parse JWK (kid=%s): %w", k.Kid, err)
		}
		keys[k.Kid] = pub
	}

	return keys, nil
}

// Validate JWT signature, issuer, and token claims
func (v *Verifier) accessTokenVerify(ctx context.Context, tokenString string) (*AccessTokenClaims, error) {
	keys, err := v.fetchJWKS(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch JWKS: %w", err)
	}

	claims := &AccessTokenClaims{}
	keyFunc := func(token *jwt.Token) (any, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("token header missing kid")
		}
		key, ok := keys[kid]
		if !ok {
			return nil, fmt.Errorf("no JWKS key found for kid %q", kid)
		}
		return key, nil
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		keyFunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Ensure token was issued for access and matches client credentials
	if claims.TokenUse != "access" {
		return nil, fmt.Errorf("expected token_use=access, got %q", claims.TokenUse)
	}
	if claims.ClientID != v.clientID {
		return nil, fmt.Errorf("token client_id mismatch")
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("token subject (sub) is missing")
	}

	return claims, nil
}

// RequireAuth enforces valid bearer access tokens and loads user identity into context
func (v *Verifier) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate Bearer Authorization header
		header := r.Header.Get("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			slog.Warn("missing or malformed Authorization header", "path", r.URL.Path)
			response.Error(w, http.StatusUnauthorized, "missing or malformed Authorization header")
			return
		}

		claims, err := v.accessTokenVerify(r.Context(), parts[1])
		if err != nil {
			slog.Warn("access token verification failed", "error", err, "path", r.URL.Path)
			if strings.Contains(err.Error(), "token is expired") {
				response.Error(w, http.StatusUnauthorized, "token_expired")
			} else {
				response.Error(w, http.StatusUnauthorized, "invalid access token")
			}
			return
		}

		// Look up user record by sub extracted from token
		sub := claims.Username
		u, err := v.userClient.GetBySub(r.Context(), sub)
		if err != nil {
			slog.Error("user resolution failed", "error", err, "sub", sub, "path", r.URL.Path)
			if user.IsNotFound(err) {
				response.Error(w, http.StatusNotFound, "user not found")
			} else {
				response.Error(w, http.StatusInternalServerError, "failed to resolve user")
			}
			return
		}

		// Store authenticated user identity into request context
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		ctx = user.WithContext(ctx, u)
		ctx = context.WithValue(ctx, userIDContextKey, u.ID)

		next(w, r.WithContext(ctx))
	}
}
