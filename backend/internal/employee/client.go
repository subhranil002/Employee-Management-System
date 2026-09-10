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

// Fetch all employees from backend
func (c *EmployeeClient) List(ctx context.Context) ([]Employee, error) {
	resp, err := c.do(ctx, http.MethodGet, "/employees", nil)
	if err != nil {
		return nil, err
	}

	slog.Info("employee client list response", "resp", resp)

	var employees []Employee
	if err := json.Unmarshal(resp.Data, &employees); err != nil {
		return nil, fmt.Errorf("decode employees list: %w", err)
	}

	return employees, nil
}

// Fetch single employee by ID from backend
func (c *EmployeeClient) Get(ctx context.Context, id string) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodGet, "/employees/"+id, nil)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode employee: %w", err)
	}

	return &emp, nil
}

// Send create employee request to backend
func (c *EmployeeClient) Create(ctx context.Context, req CreateEmployeeRequest) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodPost, "/employees", req)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode created employee: %w", err)
	}

	return &emp, nil
}

// Send partial employee update request to backend
func (c *EmployeeClient) Update(ctx context.Context, id string, req UpdateEmployeeRequest) (*Employee, error) {
	resp, err := c.do(ctx, http.MethodPatch, "/employees/"+id, req)
	if err != nil {
		return nil, err
	}

	var emp Employee
	if err := json.Unmarshal(resp.Data, &emp); err != nil {
		return nil, fmt.Errorf("decode updated employee: %w", err)
	}

	return &emp, nil
}

// Send delete employee request to backend
func (c *EmployeeClient) Delete(ctx context.Context, id string) error {
	_, err := c.do(ctx, http.MethodDelete, "/employees/"+id, nil)
	return err
}

// Execute HTTP request against backend and decode JSON envelope
func (c *EmployeeClient) do(ctx context.Context, method, path string, body any) (*crudResponse, error) {
	var reqBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&reqBody).Encode(body); err != nil {
			return nil, fmt.Errorf("encode request body: %w", err)
		}
	}

	// Build and execute HTTP request
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, &reqBody)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call employee backend: %w", err)
	}
	defer res.Body.Close()

	rawBody, _ := io.ReadAll(res.Body)
	slog.Info("raw backend response", "status", res.StatusCode, "body", string(rawBody))

	// Decode response envelope
	var envelope crudResponse
	if err := json.Unmarshal(rawBody, &envelope); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Only treat non-2xx with no data as an error
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
