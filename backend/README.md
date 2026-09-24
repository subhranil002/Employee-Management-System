# Backend API (Go Gateway & Verifier)

The Go backend serves as the secure API gateway for the Employee Management System. It validates Cognito JWT tokens, resolves authenticated user identities, and proxies tenant-scoped employee CRUD requests to the serverless Lambda microservices.

## Features

- **JWT Authentication**: Validates AWS Cognito access tokens via RS256 JWKS key sets.
- **Tenant Context Resolution**: Maps token claims (`username`/email) to MongoDB user IDs and passes them in downstream requests (`X-User-Id`).
- **REST Endpoints**:
  - `GET /healthz`: Public service health check.
  - `GET /profile`: Returns authenticated user profile.
  - `GET /employees`: Lists employees for the tenant.
  - `POST /employees`: Creates an employee under the tenant.
  - `GET /employees/{id}`: Retrieves single employee.
  - `PATCH /employees/{id}`: Partially updates employee.
  - `DELETE /employees/{id}`: Deletes employee.
- **Middleware**: Structured JSON request logging, CORS, and token verification.

## Prerequisites

- Go (1.22+)
- AWS Cognito User Pool & App Client

## Configuration

Set configuration via environment variables or AWS Secrets Manager:

```env
APP_ENV=development
HTTP_ADDR=:3000
AWS_REGION=ap-south-1
COGNITO_USER_POOL_ID=ap-south-1_xxxxxxxxx
COGNITO_CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxx
EMS_API_URL=https://<api-gateway-id>.execute-api.ap-south-1.amazonaws.com
ALLOWED_ORIGIN=http://localhost:5173
# Optional: BACKEND_SECRET_ARN=arn:aws:secretsmanager:...
```

## Running Locally

```bash
# Download dependencies
go mod download

# Run server
go run cmd/server/main.go
```
