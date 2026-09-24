# Serverless CRUD & User Functions (AWS Lambda)

Python AWS Lambda functions managing user synchronization and employee CRUD operations persisted in MongoDB.

## Functions

- **`create_user` (Cognito Post Confirmation Trigger)**:
  - Invoked automatically when a user confirms signup in AWS Cognito.
  - Upserts `cognitoSub`, `email`, and `name` into the `users` collection in MongoDB.
  - Handler: `create_user.handler.handler`.
- **`read_user`**:
  - Handles `GET /users?email=...` or `GET /users` to look up users by email.
  - Handler: `read_user.handler.handler`.
- **`create_employee`**:
  - Handles `POST /employees` scoped to the caller's `X-User-Id`.
  - Handler: `create_employee.handler.handler`.
- **`read_employee`**:
  - Handles `GET /employees` and `GET /employees/{id}` scoped to `X-User-Id`.
  - Handler: `read_employee.handler.handler`.
- **`update_employee`**:
  - Handles `PATCH /employees/{id}` scoped to `X-User-Id`.
  - Handler: `update_employee.handler.handler`.
- **`delete_employee`**:
  - Handles `DELETE /employees/{id}` scoped to `X-User-Id`.
  - Handler: `delete_employee.handler.handler`.

## Prerequisites

- Python 3.12+
- MongoDB instance (Atlas or local)

## Environment Variables

Configured in Lambda or via AWS Secrets Manager:

```env
MONGODB_URI=mongodb+srv://<user>:<password>@cluster.mongodb.net
MONGODB_DATABASE=employee_db
MONGODB_EMPLOYEE_COLLECTION=employees
MONGODB_USER_COLLECTION=users
# Optional: MONGODB_SECRET_ARN=arn:aws:secretsmanager:...
```

## Local Testing

```bash
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python app/test_local.py
```

## Packaging

Package functions along with `common/` module into zip files:

```bash
cd app
zip -r ../build/functions/create_user-v2.zip create_user/ common/
```
