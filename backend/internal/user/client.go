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
	"time"
)

type crudResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetByEmail queries the user service by email address
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

// GetBySub queries the user service by Cognito Sub
func (c *UserClient) GetBySub(ctx context.Context, sub string) (*User, error) {
	path := "/users?sub=" + url.QueryEscape(sub)

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

// Executes an HTTP request against the Lambda backend and decodes the response envelope
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

	// Handle non-2xx status codes returned by backend
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

func IsNotFound(err error) bool {
	var be *backendError
	if errors.As(err, &be) {
		return be.StatusCode == http.StatusNotFound
	}

	return false
}
