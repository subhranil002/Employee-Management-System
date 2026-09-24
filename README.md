# Employee Management System

A multi-tenant full-stack employee directory application built with React, a Go API Gateway, AWS Cognito authentication, and Python AWS Lambda serverless microservices backed by MongoDB.

## Architecture

- **[Frontend](./frontend)**: React 18 SPA (Vite + TypeScript + Tailwind CSS) with AWS Cognito OIDC authentication and silent token renewal.
- **[Backend Gateway](./backend)**: Go HTTP REST API verifying Cognito access tokens via JWKS, resolving tenant identities, and dispatching requests.
- **[Serverless Lambda Functions](./employee-crud-lambda)**: Python Lambda functions handling employee CRUD and Cognito post-confirmation user synchronization in MongoDB.

## User Lifecycle & Data Flow

1. **User Sign Up & Confirmation**:
   - The user signs up and confirms their account via AWS Cognito.
   - Cognito invokes the `create_user` Post Confirmation Lambda trigger.
   - The trigger upserts the user's `cognitoSub`, `email`, and `name` into MongoDB.
2. **Authentication & Authorization**:
   - The frontend signs in through Cognito and receives a Bearer Access Token.
   - Requests to the Go backend carry `Authorization: Bearer <access_token>`.
   - The Go gateway verifies token validity against Cognito JWKS, resolves the MongoDB user identity, and attaches `X-User-Id` to downstream Lambda calls.
3. **Employee CRUD**:
   - All employee records are scoped to `createdBy: <X-User-Id>`, ensuring strict tenant isolation.

## Component Documentation

- [Backend Documentation & Setup](./backend/README.md)
- [Serverless Functions Documentation & Setup](./employee-crud-lambda/README.md)
- [Frontend Documentation & Setup](./frontend/README.md)
