import json
import os
from urllib.parse import quote_plus

import boto3
from dotenv import load_dotenv

load_dotenv()


def get_mongodb_config():
    secret_arn = os.getenv("MONGODB_SECRET_ARN")

    # AWS: load credentials from Secrets Manager
    if secret_arn:
        secrets_client = boto3.client("secretsmanager")

        result = secrets_client.get_secret_value(
            SecretId=secret_arn
        )

        secret = json.loads(result["SecretString"])

        username = quote_plus(secret["username"])
        password = quote_plus(secret["password"])
        host = secret["host"]
        port = secret.get("port", 27017)
        database = secret.get("database", "employee_db")
        auth_source = secret.get("authSource", "admin")

        mongodb_uri = (
            f"mongodb://{username}:{password}"
            f"@{host}:{port}/{database}"
            f"?authSource={quote_plus(auth_source)}"
        )

        return {
            "uri": mongodb_uri,
            "database": database,
            "employee_collection": secret.get(
                "employee_collection",
                os.getenv(
                    "MONGODB_EMPLOYEE_COLLECTION",
                    "employees",
                ),
            ),
            "user_collection": secret.get(
                "user_collection",
                os.getenv(
                    "MONGODB_USER_COLLECTION",
                    "users",
                ),
            ),
        }

    # Local: load from .env
    mongodb_uri = os.getenv("MONGODB_URI")

    if not mongodb_uri:
        raise RuntimeError(
            "MONGODB_URI environment variable is missing "
            "and MONGODB_SECRET_ARN is not configured"
        )

    return {
        "uri": mongodb_uri,
        "database": os.getenv(
            "MONGODB_DATABASE",
            "employee_db",
        ),
        "employee_collection": os.getenv(
            "MONGODB_EMPLOYEE_COLLECTION",
            "employees",
        ),
        "user_collection": os.getenv(
            "MONGODB_USER_COLLECTION",
            "users",
        ),
    }


mongodb_config = get_mongodb_config()

MONGODB_URI = mongodb_config["uri"]
MONGODB_DATABASE = mongodb_config["database"]
employee_collection = mongodb_config["employee_collection"]
user_collection = mongodb_config["user_collection"]
