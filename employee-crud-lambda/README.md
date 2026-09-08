# Employee CRUD Serverless API

This component contains Python-based AWS Lambda functions responsible for the core Create, Read, Update, and Delete (CRUD) operations on employee data. It persists data to a MongoDB database.

## Prerequisites

- Python 3.12+
- MongoDB instance (local or Atlas)

## Setup

1. Navigate to the Lambda directory:
   ```bash
   cd employee-crud-lambda
   ```

2. Create a virtual environment and install dependencies:
   ```bash
   python -m venv .venv
   source .venv/bin/activate  # On Windows use `.venv\Scripts\activate`
   pip install -r requirements.txt
   ```

3. Configure your environment variables. Create a `.env` file in this directory (or use AWS Secrets Manager in production):

   ```env
   MONGODB_URI=mongodb://localhost:27017/
   MONGODB_DATABASE=employee_db
   MONGODB_COLLECTION=employees
   ```
   *Note: If deployed to AWS, you can use `MONGODB_SECRET_ARN` to pull database credentials securely from AWS Secrets Manager.*

## Local Testing

You can use the provided local testing script to verify the functions:

```bash
python app/test_local.py
```

## Structure

- `app/common/`: Shared configuration, database connection, and repository patterns.
- `app/create_employee/`: Lambda handler for POST requests.
- `app/read_employee/`: Lambda handler for GET requests.
- `app/update_employee/`: Lambda handler for PUT/PATCH requests.
- `app/delete_employee/`: Lambda handler for DELETE requests.

## Deployment

The functions can be packaged as ZIP files. To generate a deployment package manually:

```bash
pip install --target ./package -r requirements.txt
cd package
zip -r ../employee-crud.zip .
cd ../app
zip -g -r ../employee-crud.zip .
```

The resulting `employee-crud.zip` can be uploaded to AWS Lambda via the AWS CLI or Console.

