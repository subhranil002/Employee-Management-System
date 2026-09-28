import os

import aws_cdk as cdk
from dotenv import load_dotenv

from infrastructure.infrastructure_stack import EmployeeManagementStack

# Load environment variables from .env file
load_dotenv()

account = os.environ.get("CDK_DEFAULT_ACCOUNT")


region = os.environ.get("CDK_DEFAULT_REGION")

if not account:
    raise RuntimeError(
        "CDK_DEFAULT_ACCOUNT is not set. "
        "Check your AWS/CDK environment."
    )

if not region:
    raise RuntimeError(
        "CDK_DEFAULT_REGION is not set. "
        "Check your AWS/CDK environment."
    )


app = cdk.App()


EmployeeManagementStack(
    app,
    "EmployeeManagementStack",
    env=cdk.Environment(
        account=account,
        region=region,
    ),
)


app.synth()