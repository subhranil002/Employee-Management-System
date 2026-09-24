package user

import "context"

type User struct {
	ID    string `json:"_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userCtxKey struct{}

// WithContext attaches user identity to request context
func WithContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, u)
}

// FromContext extracts user identity from request context
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userCtxKey{}).(*User)
	return u, ok && u != nil
}
