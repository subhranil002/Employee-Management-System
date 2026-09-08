# Backend API (Auth & Gateway)

This is the Go-based backend for the Employee Management System. It primarily handles user authentication using AWS Cognito and acts as an entry point/middleware layer.

## Prerequisites

- Go (1.20+)
- AWS Cognito User Pool and Client configured

## Setup

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create your environment configuration:
   ```bash
   cp .env.example .env
   ```

4. Configure the `.env` file with your details:
   ```env
   APP_ENV=development
   HTTP_ADDR=:8080
   AWS_REGION=ap-south-1
   COGNITO_USER_POOL_ID=your-user-pool-id
   COGNITO_CLIENT_ID=your-client-id
   COGNITO_CLIENT_SECRET=your-client-secret
   ```

## Running the Server

To start the server in development mode:

```bash
go run cmd/server/main.go
```

The server will start on port `8080` (or whatever `HTTP_ADDR` is set to).

## Architecture overview

- `cmd/server/main.go`: Application entry point.
- `internal/auth/`: Contains handlers, models, and services for authentication.
- `internal/employee/`: HTTP client logic mapping to the Python AWS Lambda endpoints.
- `internal/middleware/`: CORS, authentication, token refresh, and logging middlewares.
- `pkg/`: Reusable packages for generic responses and cookie management.

