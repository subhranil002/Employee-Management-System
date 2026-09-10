package auth

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmSignUpRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserSub       string `json:"user_sub,omitempty"`
	UserConfirmed bool   `json:"user_confirmed"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int32  `json:"expires_in,omitempty"`
}

type ProfileResponse struct {
	UserSub  string `json:"user_sub"`
	Username string `json:"username"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
}

const (
	UserStatusNotFound    = "NOT_FOUND"
	UserStatusUnconfirmed = "UNCONFIRMED"
	UserStatusConfirmed   = "CONFIRMED"
)

type UserStatusResponse struct {
	Status string `json:"status"`
}
