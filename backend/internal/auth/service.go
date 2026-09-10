package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/smithy-go"
)

var (
	ErrInvalid      = errors.New("invalid request")
	ErrUnauthorized = errors.New("unauthorized")

	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

type Service struct {
	client       *cognito.Client
	userPoolID   string
	clientID     string
	clientSecret string
}

func NewService(
	client *cognito.Client,
	userPoolID string,
	clientID string,
	clientSecret string,
) *Service {
	return &Service{
		client:       client,
		userPoolID:   userPoolID,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (s *Service) Signup(
	ctx context.Context,
	req SignUpRequest,
) (*SignUpResponse, error) {

	// Validate signup input credentials
	if strings.TrimSpace(req.Name) == "" {
		return nil, fmt.Errorf("%w: name is required", ErrInvalid)
	}

	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	name := strings.TrimSpace(req.Name)

	// Build Cognito SignUp request payload
	input := &cognito.SignUpInput{
		ClientId: aws.String(s.clientID),
		Username: aws.String(email),
		Password: aws.String(req.Password),

		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
			},
		},
	}

	// Calculate secret hash when client secret is configured
	if s.clientSecret != "" {
		input.SecretHash = aws.String(
			secretHash(email, s.clientID, s.clientSecret),
		)
	}

	// Create user in Cognito User Pool
	result, err := s.client.SignUp(ctx, input)
	if err != nil {
		return nil, handleCognitoError(err)
	}

	return &SignUpResponse{
		UserSub:       aws.ToString(result.UserSub),
		UserConfirmed: result.UserConfirmed,
	}, nil
}

func (s *Service) ConfirmSignup(
	ctx context.Context,
	req ConfirmSignUpRequest,
) error {

	// Validate email and confirmation code inputs
	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("%w: confirmation code is required", ErrInvalid)
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Build Cognito ConfirmSignUp request payload
	input := &cognito.ConfirmSignUpInput{
		ClientId:         aws.String(s.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(strings.TrimSpace(req.Code)),
	}

	// Calculate secret hash when client secret is configured
	if s.clientSecret != "" {
		input.SecretHash = aws.String(
			secretHash(email, s.clientID, s.clientSecret),
		)
	}

	// Confirm registration with Cognito
	_, err := s.client.ConfirmSignUp(ctx, input)

	return handleCognitoError(err)
}

func (s *Service) Signin(
	ctx context.Context,
	req SignInRequest,
) (*SignInResponse, error) {

	// Validate signin input credentials
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	if req.Password == "" {
		return nil, fmt.Errorf("%w: password is required", ErrInvalid)
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Build auth parameters for USER_PASSWORD_AUTH flow
	authParams := map[string]string{
		"USERNAME": email,
		"PASSWORD": req.Password,
	}

	// Calculate secret hash when client secret is configured
	if s.clientSecret != "" {
		authParams["SECRET_HASH"] = secretHash(
			email,
			s.clientID,
			s.clientSecret,
		)
	}

	// Authenticate user credentials against Cognito
	input := &cognito.InitiateAuthInput{
		ClientId:       aws.String(s.clientID),
		AuthFlow:       types.AuthFlowTypeUserPasswordAuth,
		AuthParameters: authParams,
	}

	result, err := s.client.InitiateAuth(ctx, input)
	if err != nil {
		return nil, handleCognitoError(err)
	}

	if result.AuthenticationResult == nil {
		return nil, fmt.Errorf("%w: missing authentication result", ErrUnauthorized)
	}

	return &SignInResponse{
		AccessToken:  aws.ToString(result.AuthenticationResult.AccessToken),
		RefreshToken: aws.ToString(result.AuthenticationResult.RefreshToken),
		ExpiresIn:    result.AuthenticationResult.ExpiresIn,
	}, nil
}

func (s *Service) Logout(
	ctx context.Context,
	accessToken string,
) error {

	if accessToken == "" {
		return fmt.Errorf("%w: access token is required", ErrInvalid)
	}

	// Build Cognito GlobalSignOut request payload
	input := &cognito.GlobalSignOutInput{
		AccessToken: aws.String(accessToken),
	}

	// Invalidate tokens for the user in Cognito
	_, err := s.client.GlobalSignOut(ctx, input)

	return handleCognitoError(err)
}

func (s *Service) GetProfile(
	ctx context.Context,
	accessToken string,
) (*ProfileResponse, error) {

	input := &cognito.GetUserInput{
		AccessToken: aws.String(accessToken),
	}

	result, err := s.client.GetUser(ctx, input)
	if err != nil {
		return nil, handleCognitoError(err)
	}

	profile := &ProfileResponse{
		Username: aws.ToString(result.Username),
	}

	for _, attr := range result.UserAttributes {
		name := aws.ToString(attr.Name)
		val := aws.ToString(attr.Value)

		switch name {
		case "sub":
			profile.UserSub = val
		case "email":
			profile.Email = val
		case "name":
			profile.Name = val
		}
	}

	return profile, nil
}

func secretHash(
	username string,
	clientID string,
	clientSecret string,
) string {

	// Compute HMAC-SHA256 signature for client secret verification
	message := username + clientID

	mac := hmac.New(
		sha256.New,
		[]byte(clientSecret),
	)

	_, _ = mac.Write([]byte(message))

	return base64.StdEncoding.EncodeToString(
		mac.Sum(nil),
	)
}

func validateEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return fmt.Errorf("%w: email is required", ErrInvalid)
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%w: invalid email format", ErrInvalid)
	}
	return nil
}

func validatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("%w: password is required", ErrInvalid)
	}
	if len(password) < 8 {
		return fmt.Errorf("%w: password must be at least 8 characters", ErrInvalid)
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("%w: password must contain at least one uppercase letter", ErrInvalid)
	}
	if !hasLower {
		return fmt.Errorf("%w: password must contain at least one lowercase letter", ErrInvalid)
	}
	if !hasNumber {
		return fmt.Errorf("%w: password must contain at least one number", ErrInvalid)
	}
	if !hasSpecial {
		return fmt.Errorf("%w: password must contain at least one special character", ErrInvalid)
	}

	return nil
}

// Maps AWS Cognito errors to application sentinel errors for appropriate HTTP status codes
func handleCognitoError(err error) error {
	if err == nil {
		return nil
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.ErrorCode() {
		case "NotAuthorizedException", "UserNotFoundException":
			return fmt.Errorf("%w: %s", ErrUnauthorized, apiErr.ErrorMessage())
		default:
			// Map other Cognito validation/state errors to ErrInvalid (400)
			return fmt.Errorf("%w: %s", ErrInvalid, apiErr.ErrorMessage())
		}
	}
	return err
}

func (s *Service) GetUserStatus(
	ctx context.Context,
	email string,
) (string, error) {

	result, err := s.client.AdminGetUser(
		ctx,
		&cognito.AdminGetUserInput{
			UserPoolId: aws.String(s.userPoolID),
			Username:   aws.String(email),
		},
	)

	if err != nil {
		var notFound *types.UserNotFoundException

		if errors.As(err, &notFound) {
			return UserStatusNotFound, nil
		}

		return "", err
	}

	return string(result.UserStatus), nil
}
