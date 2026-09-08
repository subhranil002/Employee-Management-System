# Employee Management System

A full-stack application for managing employees, featuring a React frontend, a Go authentication backend, and a Python serverless API for CRUD operations.

## Architecture

The system is composed of three main components:

- **Frontend** (`/frontend`): A React application written in TypeScript and built with Vite. It provides the user interface for authentication and employee management.
- **Backend API Gateway / Auth** (`/backend`): A Go-based REST API that handles user authentication (via AWS Cognito), session management, and routing.
- **Serverless CRUD API** (`/employee-crud-lambda`): Python-based AWS Lambda functions that handle the core Create, Read, Update, and Delete operations for employee records, storing data in MongoDB.

## Prerequisites

- Node.js (v18+)
- Go (v1.20+)
- Python (v3.12+)
- MongoDB
- AWS Account (for Cognito and deployment)

## Getting Started

To run the project locally, you will need to start all three components. Please refer to their respective README files for detailed instructions:

1. [Backend (Go API)](./backend/README.md)
2. [Serverless CRUD API (Python)](./employee-crud-lambda/README.md)
3. [Frontend (React)](./frontend/README.md)

## Environment Variables

Each component requires specific environment variables. Make sure to configure the `.env` files in each respective directory based on the `.env.example` or README instructions provided.

## Deployment

This repository is structured as a monorepo and is ready to be pushed to GitHub.
- The **frontend** can be deployed to Vercel, Netlify, or AWS S3/CloudFront.
- The **backend** can be containerized or deployed to AWS ECS, App Runner, or any standard VM/PaaS.
- The **employee-crud-lambda** service is designed to be deployed to AWS Lambda using AWS SAM, Serverless Framework, or AWS CDK.

