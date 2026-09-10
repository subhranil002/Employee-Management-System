package middleware

import (
	"context"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	TokenUse string `json:"token_use"`
	ClientID string `json:"client_id"`
	Scope    string `json:"scope"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}

type TokenVerifier struct {
	jwks      keyfunc.Keyfunc
	issuer    string
	clientID  string
	refresher *TokenRefresher
}

func NewTokenVerifier(
	ctx context.Context,
	issuer string,
	clientID string,
	refresher *TokenRefresher,
) (*TokenVerifier, error) {

	jwksURL := issuer + "/.well-known/jwks.json"

	// Fetch and cache Cognito JWKS public keys for token signature verification
	jwks, err := keyfunc.NewDefaultCtx(
		ctx,
		[]string{jwksURL},
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create cognito JWKS verifier: %w",
			err,
		)
	}

	return &TokenVerifier{
		jwks:      jwks,
		issuer:    issuer,
		clientID:  clientID,
		refresher: refresher,
	}, nil
}

func (v *TokenVerifier) Verify(
	tokenString string,
) (*AccessTokenClaims, error) {

	claims := &AccessTokenClaims{}

	// Parse and validate RS256 token against Cognito JWKS public keys
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		v.jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf("verify cognito token: %w", err)
	}

	// Validate standard Cognito claims
	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	if claims.TokenUse != "access" {
		return nil, fmt.Errorf("token is not an access token")
	}

	if claims.ClientID != v.clientID {
		return nil, fmt.Errorf("token belongs to another app client")
	}

	if claims.Subject == "" {
		return nil, fmt.Errorf("token subject is missing")
	}

	return claims, nil
}

// Decode token claims without signature or expiration check
func parseExpiredClaims(tokenString string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}
