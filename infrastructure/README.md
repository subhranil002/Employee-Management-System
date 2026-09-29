# Infrastructure (AWS CDK with Python)

Infrastructure-as-Code (IaC) for the Employee Management System, built using [AWS Cloud Development Kit (AWS CDK v2)](https://docs.aws.amazon.com/cdk/v2/guide/home.html) and Python 3.12.

This stack provisions the complete serverless backend foundation: AWS Lambda functions, a shared Python dependency layer, an AWS API Gateway HTTP API, and IAM role permissions for AWS Secrets Manager access.

---

## Architecture Overview

```mermaid
flowchart TD
    subgraph AWS Cloud
        subgraph Secrets Manager
            SM[MongoDB Secret<br/>MONGODB_SECRET_ARN]
        end

        subgraph Lambda Infrastructure
            Layer[Shared Python Layer<br/>employee-management-python-dependencies]
            CU[create_user Lambda<br/>Cognito Trigger]
            RU[read_user Lambda]
            CE[create_employee Lambda]
            RE[read_employee Lambda]
            UE[update_employee Lambda]
            DE[delete_employee Lambda]
        end

        subgraph API Gateway
            APIGW[HTTP API v2<br/>employee-management-api]
        end

        subgraph Cognito
            COG[Cognito User Pool]
        end
    end

    SM -.->|Grant Read| Layer
    SM -.->|Grant Read| CU & RU & CE & RE & UE & DE
    Layer -.->|Mounted to| CU & RU & CE & RE & UE & DE

    COG -->|Post-Confirmation Trigger| CU
    APIGW -->|GET /users| RU
    APIGW -->|POST /employees| CE
    APIGW -->|GET /employees<br/>GET /employees/{id}| RE
    APIGW -->|PATCH /employees/{id}| UE
    APIGW -->|DELETE /employees/{id}| DE
```

---

## Provisioned Resources

### 1. Shared Python Dependency Layer
- **Construct**: `aws_lambda.LayerVersion`
- **Name**: `employee-management-python-dependencies`
- **Runtime**: Python 3.12 (`x86_64`)
- **Source**: `employee-crud-lambda/requirements.txt`
- **Bundling**: Built using Docker (`public.ecr.aws/sam/build-python3.12`) into `/asset-output/python`.
- **Asset Hashing**: Uses custom deterministic SHA-256 hashing so the layer is only rebuilt when `requirements.txt` changes.

### 2. AWS Lambda Functions
Six Python 3.12 functions configured with 512 MB memory, 30s timeout, and the shared dependency layer:

| Function Name | Handler | Purpose | Integration |
|---|---|---|---|
| `create_user` | `create_user.handler.handler` | Synchronizes confirmed Cognito users to MongoDB | Cognito Post-Confirmation Trigger |
| `read_user` | `read_user.handler.handler` | Queries user by Cognito `sub` or `email` | HTTP API (`GET /users`) |
| `create_employee` | `create_employee.handler.handler` | Creates employee scoped to `X-User-Id` | HTTP API (`POST /employees`) |
| `read_employee` | `read_employee.handler.handler` | Lists or gets employee scoped to `X-User-Id` | HTTP API (`GET /employees`, `GET /employees/{id}`) |
| `update_employee` | `update_employee.handler.handler` | Updates employee scoped to `X-User-Id` | HTTP API (`PATCH /employees/{id}`) |
| `delete_employee` | `delete_employee.handler.handler` | Deletes employee scoped to `X-User-Id` | HTTP API (`DELETE /employees/{id}`) |

> **Asset Hashing**: Each function's CDK asset hash is computed only from `app/<function_name>/` and `app/common/`. Changes in one Lambda handler will not force rebuilds of unrelated functions.

### 3. Amazon API Gateway HTTP API (v2)
- **API Name**: `employee-management-api`
- **Protocol**: HTTP
- **Routes & Integrations**:
  - `GET /users` &rarr; `read_user` Lambda
  - `POST /employees` &rarr; `create_employee` Lambda
  - `GET /employees` &rarr; `read_employee` Lambda
  - `GET /employees/{id}` &rarr; `read_employee` Lambda
  - `PATCH /employees/{id}` &rarr; `update_employee` Lambda
  - `DELETE /employees/{id}` &rarr; `delete_employee` Lambda

### 4. AWS Secrets Manager Integration
- References existing MongoDB Secret ARN via `MONGODB_SECRET_ARN`.
- Grants read permissions (`mongodb_secret.grant_read(function)`) to each Lambda execution role.
- Passes `MONGODB_SECRET_ARN` into Lambda environment variables.

---

## Stack Outputs

Upon successful deployment, the stack exports:

| Output Key | Description | Downstream Consumer |
|---|---|---|
| `HttpApiUrl` | Endpoint URL of the API Gateway HTTP API | Go Backend (`EMS_API_URL`) |
| `CreateUserLambdaArn` | ARN of the `create_user` Lambda | AWS Cognito User Pool Trigger |
| `ReadUserLambdaArn` | ARN of the `read_user` Lambda | Internal inspection |
| `CreateEmployeeLambdaArn` | ARN of the `create_employee` Lambda | Internal inspection |
| `ReadEmployeeLambdaArn` | ARN of the `read_employee` Lambda | Internal inspection |
| `UpdateEmployeeLambdaArn` | ARN of the `update_employee` Lambda | Internal inspection |
| `DeleteEmployeeLambdaArn` | ARN of the `delete_employee` Lambda | Internal inspection |
| `DependencyLayerArn` | ARN of the published Lambda Layer | Internal inspection |

---

## Prerequisites

1. **Python 3.12+**
2. **Node.js 18+** & **AWS CDK CLI**:
   ```bash
   npm install -g aws-cdk
   cdk --version
   ```
3. **Docker Desktop / Docker Daemon**: Running locally (required for CDK Lambda layer and function asset bundling).
4. **AWS CLI v2**: Configured with credentials for your AWS account (`aws configure`).

---

## Configuration & Environment Variables

Create a `.env` file inside the `infrastructure/` directory:

```env
CDK_DEFAULT_ACCOUNT=123456789012
CDK_DEFAULT_REGION=ap-south-1
MONGODB_SECRET_ARN=arn:aws:secretsmanager:ap-south-1:123456789012:secret:MONGODB_SECRET_ARN-xxxxxx
```

| Variable | Required | Description |
|---|---|---|
| `CDK_DEFAULT_ACCOUNT` | Yes | AWS Account ID where resources are deployed |
| `CDK_DEFAULT_REGION` | Yes | AWS Region (e.g., `ap-south-1`, `us-east-1`) |
| `MONGODB_SECRET_ARN` | Optional | ARN of the Secrets Manager secret storing MongoDB connection details |

---

## Deployment Guide

### 1. Initialize Virtual Environment
```bash
cd infrastructure

# On Linux / macOS:
python3 -m venv .venv
source .venv/bin/activate

# On Windows:
.venv\Scripts\activate.bat
```

### 2. Install CDK Dependencies
```bash
pip install -r requirements.txt
```

### 3. Bootstrap CDK (First Time Only)
If you have not bootstrapped the target AWS account and region for CDK:
```bash
cdk bootstrap aws://$CDK_DEFAULT_ACCOUNT/$CDK_DEFAULT_REGION
```

### 4. Synthesize CloudFormation Template
```bash
cdk synth
```

### 5. Inspect Differences
```bash
cdk diff
```

### 6. Deploy the Stack
```bash
cdk deploy
```

---

## Post-Deployment Wiring

1. **Connect Cognito Post-Confirmation Trigger**:
   - Navigate to the **AWS Cognito Console** &rarr; Select your User Pool.
   - Go to **User pool properties** &rarr; **Lambda triggers**.
   - Under **Sign-up**, select **Post confirmation trigger** &rarr; Select the Lambda function corresponding to `CreateUserLambdaArn`.
   - Grant Cognito permission to invoke the Lambda:
     ```bash
     aws lambda add-permission \
       --function-name <CreateUserLambdaArn> \
       --statement-id CognitoPostConfirmation \
       --action lambda:InvokeFunction \
       --principal cognito-idp.amazonaws.com \
       --source-arn arn:aws:cognito-idp:<region>:<account>:userpool/<user-pool-id>
     ```

2. **Connect Go Backend Gateway**:
   - Copy `HttpApiUrl` from the CDK deployment outputs.
   - Set `EMS_API_URL=<HttpApiUrl>` in `backend/.env`.
