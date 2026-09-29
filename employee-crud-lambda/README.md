# Serverless CRUD & User Functions (AWS Lambda)

Python 3.12 AWS Lambda microservices for tenant user synchronization and multi-tenant employee CRUD operations persisted in MongoDB.

---

## Architecture Overview

```mermaid
flowchart TD
    subgraph Triggers & Ingress
        COG[Cognito Post-Confirmation Trigger]
        APIGW[API Gateway HTTP API v2]
    end

    subgraph Lambda Microservices
        CU[create_user]
        RU[read_user]
        CE[create_employee]
        RE[read_employee]
        UE[update_employee]
        DE[delete_employee]
    end

    subgraph Common Module
        SVC[service.py<br/>Business Logic & Validation]
        SCH[schemas.py<br/>Pydantic Models]
        REPO[repository.py<br/>PyMongo Repositories]
        DB[database.py & config.py<br/>Connection & Secrets]
        RESP[response.py<br/>API Response Helper]
    end

    subgraph Storage
        SM[AWS Secrets Manager]
        MONGO[(MongoDB / Atlas / DocumentDB)]
    end

    COG -->|User confirmed| CU
    APIGW -->|GET /users| RU
    APIGW -->|POST /employees| CE
    APIGW -->|GET /employees<br/>GET /employees/{id}| RE
    APIGW -->|PATCH /employees/{id}| UE
    APIGW -->|DELETE /employees/{id}| DE

    CU & RU & CE & RE & UE & DE --> Common Module
    DB -.->|Fetch credentials| SM
    REPO -->|Query / Mutate| MONGO
```

---

## Function Reference

### 1. `create_user` (Cognito Post-Confirmation Trigger)
- **Handler**: `create_user.handler.handler`
- **Trigger**: AWS Cognito Post-Confirmation Trigger
- **Flow**:
  1. Triggered immediately after a user verifies their sign-up in AWS Cognito.
  2. Receives `event["request"]["userAttributes"]` containing `sub`, `email`, and optionally `name`, `given_name`, `family_name`.
  3. Upserts the user record into the `users` collection in MongoDB keyed on `cognitoSub`.
  4. Returns the unmodified event object to AWS Cognito.

### 2. `read_user`
- **Handler**: `read_user.handler.handler`
- **Method / Path**: `GET /users`
- **Query Parameters**:
  - `?sub=<cognitoSub>`: Look up user by Cognito Subject identifier (used by Go gateway).
  - `?email=<email>`: Look up user by email address.
  - (No params): Returns a list of all users.
- **Success Response** (`200 OK`):
  ```json
  {
    "success": true,
    "message": "User retrieved successfully",
    "data": {
      "_id": "64d0a1b2c3d4e5f6a7b8c9d0",
      "cognitoSub": "7a35c102-...",
      "name": "Jane Doe",
      "email": "jane@example.com"
    }
  }
  ```

### 3. `create_employee`
- **Handler**: `create_employee.handler.handler`
- **Method / Path**: `POST /employees`
- **Required Headers**: `X-User-Id: <MongoDB User ID>` (injected by Go gateway)
- **Body**:
  ```json
  {
    "empID": "EMP-001",
    "name": "Alice Smith",
    "email": "alice@example.com",
    "department": "Engineering",
    "salary": 95000.0
  }
  ```
- **Validation Rules**:
  - `empID` must be unique per tenant (`createdBy`).
  - `email` must be valid and unique per tenant (`createdBy`).
  - `salary` must be &ge; 0.
- **Success Response** (`201 Created`):
  ```json
  {
    "success": true,
    "message": "Employee created successfully",
    "data": {
      "_id": "67890abcdef...",
      "empID": "EMP-001",
      "createdBy": "64d0a1b2c3d4e5f6a7b8c9d0",
      "name": "Alice Smith",
      "email": "alice@example.com",
      "department": "Engineering",
      "salary": 95000.0
    }
  }
  ```

### 4. `read_employee`
- **Handler**: `read_employee.handler.handler`
- **Method / Path**:
  - `GET /employees`: Returns all employees where `createdBy == X-User-Id`, sorted by `name`.
  - `GET /employees/{id}`: Returns single employee matching both `_id == id` and `createdBy == X-User-Id`.
- **Required Headers**: `X-User-Id: <MongoDB User ID>`

### 5. `update_employee`
- **Handler**: `update_employee.handler.handler`
- **Method / Path**: `PATCH /employees/{id}`
- **Required Headers**: `X-User-Id: <MongoDB User ID>`
- **Body** (partial fields):
  ```json
  {
    "name": "Alice Johnson",
    "salary": 105000.0
  }
  ```
- **Validation**: Checks for unique email constraint within the tenant if email is updated.

### 6. `delete_employee`
- **Handler**: `delete_employee.handler.handler`
- **Method / Path**: `DELETE /employees/{id}`
- **Required Headers**: `X-User-Id: <MongoDB User ID>`
- **Success Response** (`200 OK`):
  ```json
  {
    "success": true,
    "message": "Employee deleted successfully",
    "data": null
  }
  ```

---

## Common Module Architecture

The `app/common` directory contains shared utilities used across all Lambda handlers:

- **`config.py`**: Resolves MongoDB configuration. When `MONGODB_SECRET_ARN` is provided, fetches credentials from AWS Secrets Manager using Boto3; otherwise falls back to local `MONGODB_URI` environment variable.
- **`database.py`**: Lazily establishes and caches a singleton `MongoClient` instance with connection timeouts and ping verification.
- **`schemas.py`**: Pydantic models with field constraints:
  - `EmployeeCreate`: `empID`, `createdBy`, `name`, `email`, `department`, `salary`.
  - `EmployeeUpdate`: Optional partial update fields.
  - `CognitoUserCreate`: `cognitoSub`, `name`, `email`.
- **`repository.py`**: Encapsulates PyMongo queries:
  - `EmployeeRepository`: Scopes all queries (`find_all`, `find_by_id`, `insert`, `update_by_id`, `delete_by_id`) with `createdBy` to ensure tenant isolation. Converts `ObjectId` to string.
  - `UserRepository`: Handles `find_by_cognito_sub`, `find_by_email`, and `upsert_by_cognito_sub`.
- **`service.py`**: Implements business rules (duplicate email/empID checks, error classification) and throws typed exceptions (`EmployeeNotFoundError`, `InvalidEmployeeError`, etc.).
- **`response.py`**: Formats API Gateway proxy responses with status code, standard JSON envelope, and CORS headers.

---

## Database Collections & Schemas

### `users` Collection
```json
{
  "_id": ObjectId("64d0a1b2c3d4e5f6a7b8c9d0"),
  "cognitoSub": "7a35c102-4567-4890-a1b2-c3d4e5f6a7b8",
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

### `employees` Collection
```json
{
  "_id": ObjectId("67890abcdef1234567890abc"),
  "empID": "EMP-001",
  "createdBy": "64d0a1b2c3d4e5f6a7b8c9d0",
  "name": "Alice Smith",
  "email": "alice@example.com",
  "department": "Engineering",
  "salary": 95000.0
}
```

> **Multi-Tenancy Model**: Every employee document contains `createdBy` storing the string representation of the tenant's MongoDB user `_id`. Queries always enforce `{"_id": ObjectId(id), "createdBy": created_by}`.

---

## Configuration & Environment Variables

When running locally or outside AWS Lambda, set the following environment variables:

```env
MONGODB_URI=mongodb+srv://<user>:<password>@cluster.mongodb.net/employee_db?retryWrites=true&w=majority
MONGODB_DATABASE=employee_db
MONGODB_EMPLOYEE_COLLECTION=employees
MONGODB_USER_COLLECTION=users

# Optional: When running in AWS Lambda
# MONGODB_SECRET_ARN=arn:aws:secretsmanager:ap-south-1:123456789012:secret:MONGODB_SECRET_ARN-xxxxxx
```

### Secrets Manager JSON Structure (if using `MONGODB_SECRET_ARN`)
```json
{
  "username": "<db_user>",
  "password": "<db_password>",
  "host": "<cluster_host>",
  "port": 27017,
  "database": "employee_db",
  "authSource": "admin"
}
```

---

## Local Development & Testing

A standalone test script `app/test_local.py` is included to verify database connectivity and basic repository operations:

```bash
cd employee-crud-lambda

# 1. Create and activate virtual environment
python3 -m venv .venv
source .venv/bin/activate

# 2. Install dependencies
pip install -r requirements.txt

# 3. Export environment variables
export MONGODB_URI="mongodb://localhost:27017/employee_db"

# 4. Run local test script
python app/test_local.py
```

---

## Deployment & Packaging

- **Automatic (Recommended)**: Use the [AWS CDK Stack](../infrastructure/README.md), which automatically bundles `requirements.txt` into a shared layer and packages functions with their `common/` dependencies.
- **Manual Zip Packaging**:
  ```bash
  cd app
  zip -r ../build/functions/create_user-v2.zip create_user/ common/
  zip -r ../build/functions/read_user-v2.zip read_user/ common/
  zip -r ../build/functions/create_employee-v2.zip create_employee/ common/
  zip -r ../build/functions/read_employee-v2.zip read_employee/ common/
  zip -r ../build/functions/update_employee-v2.zip update_employee/ common/
  zip -r ../build/functions/delete_employee-v2.zip delete_employee/ common/
  ```
