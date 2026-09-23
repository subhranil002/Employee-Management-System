# Frontend UI (React + Vite + TypeScript)

Single-page application for the Employee Management System, built with React, Vite, Tailwind CSS, and `react-oidc-context`.

## Features

- **AWS Cognito Authentication**: Authorization Code flow with PKCE via `react-oidc-context`.
- **Silent Token Refresh**: Axios interceptors automatically renew expired tokens using Cognito silent sign-in.
- **Tenant Isolation**: Only displays and manages employees created by the authenticated user.
- **Employee Management**: Filterable table with full CRUD (Add, Edit, View, Delete) modals and toast alerts.

## Prerequisites

- Node.js (v18+)
- npm or yarn

## Environment Variables

Create `.env` in the `frontend` root:

```env
VITE_API_BASE_URL=http://localhost:3000
VITE_COGNITO_AUTHORITY=https://cognito-idp.<region>.amazonaws.com/<user-pool-id>
VITE_COGNITO_CLIENT_ID=<cognito-client-id>
VITE_COGNITO_REDIRECT_URI=http://localhost:5173
VITE_COGNITO_DOMAIN=https://<your-auth-domain>.auth.<region>.amazoncognito.com
VITE_COGNITO_SCOPE=openid email phone
```

## Getting Started

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```
