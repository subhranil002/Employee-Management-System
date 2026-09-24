package employee

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type crudResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type EmployeeClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *EmployeeClient {
	return &EmployeeClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// List retrieves all employees belonging to the specified user
func (c *EmployeeClient) List(ctx context.Context, userID string) ([]Employee, error) {
	resp, err := c.do(ctx, http.MethodGet, "/employees", userID, nil)
	if err != nil {
		return nil, err
	}

	var employees []Employee
	if err := json.Unmarshal(resp.Data, &employees); err != nil {
		return nil, fmt.Errorf("decode employees list: %w", err)
	}

	return employees, nil
}

// Get fetches a single employee record by ID scoped to the user
func (c *EmployeeClient) Get(ctx context.Context, userID, id string) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodGet, "/employees/"+id, userID, nil)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode employee: %w", err)
	}

	return &emp, nil
}

// Create inserts a new employee record under the user
func (c *EmployeeClient) Create(ctx context.Context, userID string, req CreateEmployeeRequest) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodPost, "/employees", userID, req)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode created employee: %w", err)
	}

	return &emp, nil
}

// Update modifies an existing employee record belonging to the user
func (c *EmployeeClient) Update(ctx context.Context, userID, id string, req UpdateEmployeeRequest) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodPatch, "/employees/"+id, userID, req)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode updated employee: %w", err)
	}

	return &emp, nil
}

// Delete removes an employee record belonging to the user
func (c *EmployeeClient) Delete(ctx context.Context, userID, id string) error {
	_, err := c.do(ctx, http.MethodDelete, "/employees/"+id, userID, nil)
	return err
}

// Executes an HTTP request against the employee Lambda endpoint
func (c *EmployeeClient) do(ctx context.Context, method, path, userID string, body any) (*crudResponse, error) {
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

	// Propagate tenant user ID header to downstream Lambda
	if userID != "" {
		req.Header.Set("X-User-Id", userID)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call employee backend: %w", err)
	}
	defer res.Body.Close()

	rawBody, _ := io.ReadAll(res.Body)
	slog.Info("raw backend response", "status", res.StatusCode, "body", string(rawBody))

	var envelope crudResponse
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Return error for non-2xx status responses
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
