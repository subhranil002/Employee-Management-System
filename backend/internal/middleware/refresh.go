package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type RefreshedTokens struct {
	AccessToken string
	ExpiresIn   int32
}

type TokenRefresher struct {
	client       *cognito.Client
	clientID     string
	clientSecret string
}

func NewTokenRefresher(
	client *cognito.Client,
	clientID string,
	clientSecret string,
) *TokenRefresher {
	return &TokenRefresher{
		client:       client,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// Exchange refresh token for new access token with Cognito
func (r *TokenRefresher) Refresh(
	ctx context.Context,
	refreshToken string,
	username string,
) (*RefreshedTokens, error) {

	authParams := map[string]string{
		"REFRESH_TOKEN": refreshToken,
	}

	// Compute client secret hash when app client secret is configured
	if r.clientSecret != "" {
		authParams["SECRET_HASH"] = refreshSecretHash(
			username,
			r.clientID,
			r.clientSecret,
		)
	}

	// Request new tokens using Cognito REFRESH_TOKEN_AUTH flow
	result, err := r.client.InitiateAuth(
		ctx,
		&cognito.InitiateAuthInput{
			AuthFlow:       types.AuthFlowTypeRefreshTokenAuth,
			ClientId:       aws.String(r.clientID),
			AuthParameters: authParams,
		},
	)
	if err != nil {
		return nil, err
	}

	authResult := result.AuthenticationResult

	return &RefreshedTokens{
		AccessToken: aws.ToString(authResult.AccessToken),
		ExpiresIn:   authResult.ExpiresIn,
	}, nil
}

// Compute HMAC-SHA256 secret hash for Cognito authentication
func refreshSecretHash(username, clientID, clientSecret string) string {
	message := username + clientID

	mac := hmac.New(sha256.New, []byte(clientSecret))
	_, _ = mac.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
