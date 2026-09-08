# Frontend UI

This is the React frontend for the Employee Management System, built with Vite and TypeScript. It interacts with the backend Go API for authentication and the Lambda functions for employee operations.

## Prerequisites

- Node.js (v18+)
- npm or yarn or pnpm

## Setup

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Create a `.env` file in the root of the frontend directory:
   ```env
   VITE_API_BASE_URL=http://localhost:8080
   ```

## Development

To start the development server with Hot Module Replacement (HMR):

```bash
npm run dev
```

The app will typically be available at `http://localhost:5173`.

## Building for Production

To create a production-ready build:

```bash
npm run build
```

The compiled static assets will be output to the `dist/` directory, which can be deployed to static hosting providers like Vercel, Netlify, or AWS S3.

## Linting

To run ESLint and check for code quality issues:

```bash
npm run lint
```
