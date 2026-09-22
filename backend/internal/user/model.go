package user

import "context"

// User represents a user record from the backend
type User struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUserRequest is the payload sent to POST /users
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userCtxKey struct{}

// WithContext stores User in context
func WithContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

// FromContext retrieves User from context
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userCtxKey{}).(*User)
	return u, ok && u != nil
}
