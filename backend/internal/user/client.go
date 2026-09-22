package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type crudResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// UserClient is an HTTP client for the Users Lambda API
type UserClient struct {
	baseURL    string
	httpClient *http.Client

	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		locks: make(map[string]*sync.Mutex),
	}
}

// GetByEmail fetches a single user by email address
func (c *UserClient) GetByEmail(ctx context.Context, email string) (*User, error) {
	path := "/users?email=" + url.QueryEscape(email)

	resp, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(resp.Data, &u); err != nil {
		return nil, fmt.Errorf("decode user: %w", err)
	}

	return &u, nil
}

// Create creates a new user record in the backend
func (c *UserClient) Create(ctx context.Context, name, email string) (*User, error) {
	resp, err := c.do(ctx, http.MethodPost, "/users", CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(resp.Data, &u); err != nil {
		return nil, fmt.Errorf("decode created user: %w", err)
	}

	return &u, nil
}

// getLock returns a dedicated mutex for an email
func (c *UserClient) getLock(email string) *sync.Mutex {
	c.mu.Lock()
	defer c.mu.Unlock()

	l, ok := c.locks[email]
	if !ok {
		l = &sync.Mutex{}
		c.locks[email] = l
	}
	return l
}

// GetOrCreate fetches a user by email, auto-creating them if they do not exist
func (c *UserClient) GetOrCreate(ctx context.Context, name, email string) (*User, error) {
	email = strings.ToLower(strings.TrimSpace(email))

	l := c.getLock(email)
	l.Lock()
	defer l.Unlock()

	u, err := c.GetByEmail(ctx, email)
	if err == nil {
		return u, nil
	}

	if IsNotFound(err) {
		slog.Info("user not found, auto-creating", "email", email)
		return c.Create(ctx, name, email)
	}

	return nil, fmt.Errorf("get user by email: %w", err)
}

// do executes an HTTP request against the Users API and decodes the JSON envelope
func (c *UserClient) do(ctx context.Context, method, path string, body any) (*crudResponse, error) {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call user backend: %w", err)
	}
	defer res.Body.Close()

	rawBody, _ := io.ReadAll(res.Body)
	slog.Info("raw user backend response", "status", res.StatusCode, "body", string(rawBody))

	var envelope crudResponse
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, &backendError{
			StatusCode: res.StatusCode,
			Message:    envelope.Message,
		}
	}

	return &envelope, nil
}

type backendError struct {
	StatusCode int
	Message    string
}

func (e *backendError) Error() string {
	return e.Message
}

// IsNotFound returns true if the error represents a 404 from the backend
func IsNotFound(err error) bool {
	var be *backendError
	if errors.As(err, &be) {
		return be.StatusCode == http.StatusNotFound
	}

	return false
}
