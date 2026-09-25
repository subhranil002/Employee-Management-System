import os
import hashlib
from pathlib import Path

import aws_cdk as cdk
from aws_cdk import (
    BundlingOptions,
    Duration,
    aws_apigatewayv2 as apigwv2,
    aws_apigatewayv2_integrations as integrations,
    aws_lambda as lambda_,
    aws_secretsmanager as secretsmanager,
)
from constructs import Construct


# =============================================================
# PROJECT PATHS
# =============================================================

PROJECT_ROOT = Path(__file__).resolve().parents[2]

LAMBDA_ROOT = (
    PROJECT_ROOT / "employee-crud-lambda"
)

APP_ROOT = (
    LAMBDA_ROOT / "app"
)


# =============================================================
# MONGODB SECRET ARN
# =============================================================

MONGODB_SECRET_ARN = os.environ.get(
    "MONGODB_SECRET_ARN",
    "arn:aws:secretsmanager:us-east-1:121141733713:secret:MONGODB_SECRET_ARN-aY3sHA",
)



# =============================================================
# HASH HELPER
# =============================================================

def calculate_hash(*paths: Path) -> str:
    """
    Generate a deterministic SHA-256 hash from the supplied
    files/directories.

    The hash includes:
      - relative file path
      - file contents

    This allows CDK assets to change only when the files that
    actually belong to that asset change.
    """

    sha = hashlib.sha256()

    files: list[Path] = []

    for path in paths:

        if path.is_file():
            files.append(path)

        elif path.is_dir():
            files.extend(
                file
                for file in path.rglob("*")
                if file.is_file()
            )

    for file_path in sorted(files):

        relative_path = file_path.relative_to(
            PROJECT_ROOT
        )

        sha.update(
            str(relative_path).encode("utf-8")
        )

        sha.update(
            file_path.read_bytes()
        )

    return sha.hexdigest()


# =============================================================
# STACK
# =============================================================

class EmployeeManagementStack(cdk.Stack):

    def __init__(
        self,
        scope: Construct,
        construct_id: str,
        **kwargs,
    ) -> None:

        super().__init__(
            scope,
            construct_id,
            **kwargs,
        )

        # =========================================================
        # EXISTING MONGODB SECRET
        # =========================================================

        mongodb_secret = (
            secretsmanager.Secret.from_secret_complete_arn(
                self,
                "MongoDbSecret",
                MONGODB_SECRET_ARN,
            )
        )

        # =========================================================
        # SHARED DEPENDENCY LAYER
        # =========================================================

        # IMPORTANT:
        # Only requirements.txt contributes to the layer hash.

        layer_hash = calculate_hash(
            LAMBDA_ROOT / "requirements.txt"
        )

        dependency_layer = (
            lambda_.LayerVersion(
                self,
                "EmployeeManagementPythonLayer",

                layer_version_name=(
                    "employee-management-python-dependencies"
                ),

                description=(
                    "Shared Python dependencies for "
                    "Employee Management Lambda functions"
                ),

                compatible_runtimes=[
                    lambda_.Runtime.PYTHON_3_12,
                ],

                compatible_architectures=[
                    lambda_.Architecture.X86_64,
                ],

                code=cdk.aws_lambda.Code.from_asset(
                    str(LAMBDA_ROOT),

                    asset_hash=layer_hash,

                    bundling=BundlingOptions(
                        image=(
                            lambda_.Runtime.PYTHON_3_12
                            .bundling_image
                        ),

                        command=[
                            "bash",
                            "-c",

                            (
                                "mkdir -p "
                                "/asset-output/python "
                                "&& "
                                "pip install "
                                "--no-cache-dir "
                                "-r "
                                "/asset-input/requirements.txt "
                                "-t "
                                "/asset-output/python"
                            ),
                        ],
                    ),
                ),

                removal_policy=(
                    cdk.RemovalPolicy.RETAIN
                ),
            )
        )

        # =========================================================
        # LAMBDA CREATION HELPER
        # =========================================================

        def create_lambda(
            construct_id: str,
            function_name: str,
        ) -> lambda_.Function:

            # -----------------------------------------------------
            # IMPORTANT:
            #
            # This Lambda's hash depends ONLY on:
            #
            #   app/<function_name>/
            #   app/common/
            #
            # Therefore a change to another Lambda does not
            # change this function's asset hash.
            # -----------------------------------------------------

            function_hash = calculate_hash(
                APP_ROOT / function_name,
                APP_ROOT / "common",
            )

            code = (
                cdk.aws_lambda.Code.from_asset(
                    str(LAMBDA_ROOT),

                    asset_hash=function_hash,

                    bundling=BundlingOptions(
                        image=(
                            lambda_.Runtime.PYTHON_3_12
                            .bundling_image
                        ),

                        command=[
                            "bash",
                            "-c",

                            (
                                "mkdir -p /asset-output "
                                "&& "
                                f"cp -R "
                                f"/asset-input/app/"
                                f"{function_name} "
                                "/asset-output/ "
                                "&& "
                                "cp -R "
                                "/asset-input/app/common "
                                "/asset-output/"
                            ),
                        ],
                    ),
                )
            )

            function = lambda_.Function(
                self,
                construct_id,

                function_name=function_name,

                runtime=(
                    lambda_.Runtime.PYTHON_3_12
                ),

                handler=(
                    f"{function_name}.handler.handler"
                ),

                architecture=(
                    lambda_.Architecture.X86_64
                ),

                code=code,

                layers=[
                    dependency_layer,
                ],

                timeout=(
                    Duration.seconds(30)
                ),

                memory_size=512,

                environment={
                    "MONGODB_SECRET_ARN": (
                        MONGODB_SECRET_ARN
                    ),
                },
            )

            # -----------------------------------------------------
            # Allow Lambda to read the existing secret
            # -----------------------------------------------------

            mongodb_secret.grant_read(
                function
            )

            return function

        # =========================================================
        # SIX LAMBDAS
        # =========================================================

        create_user = create_lambda(
            "CreateUserFunction",
            "create_user",
        )

        read_user = create_lambda(
            "ReadUserFunction",
            "read_user",
        )

        create_employee = create_lambda(
            "CreateEmployeeFunction",
            "create_employee",
        )

        read_employee = create_lambda(
            "ReadEmployeeFunction",
            "read_employee",
        )

        update_employee = create_lambda(
            "UpdateEmployeeFunction",
            "update_employee",
        )

        delete_employee = create_lambda(
            "DeleteEmployeeFunction",
            "delete_employee",
        )

        # =========================================================
        # ONE HTTP API
        # =========================================================

        http_api = apigwv2.HttpApi(
            self,
            "EmployeeHttpApi",

            api_name=(
                "employee-management-api"
            ),
        )

        # =========================================================
        # INTEGRATIONS
        # =========================================================

        read_user_integration = (
            integrations.HttpLambdaIntegration(
                "ReadUserIntegration",
                read_user,
            )
        )

        create_employee_integration = (
            integrations.HttpLambdaIntegration(
                "CreateEmployeeIntegration",
                create_employee,
            )
        )

        read_employee_integration = (
            integrations.HttpLambdaIntegration(
                "ReadEmployeeIntegration",
                read_employee,
            )
        )

        update_employee_integration = (
            integrations.HttpLambdaIntegration(
                "UpdateEmployeeIntegration",
                update_employee,
            )
        )

        delete_employee_integration = (
            integrations.HttpLambdaIntegration(
                "DeleteEmployeeIntegration",
                delete_employee,
            )
        )

        # =========================================================
        # GET /users
        # =========================================================

        http_api.add_routes(
            path="/users",

            methods=[
                apigwv2.HttpMethod.GET,
            ],

            integration=(
                read_user_integration
            ),
        )

        # =========================================================
        # POST /employees
        # =========================================================

        http_api.add_routes(
            path="/employees",

            methods=[
                apigwv2.HttpMethod.POST,
            ],

            integration=(
                create_employee_integration
            ),
        )

        # =========================================================
        # GET /employees
        # =========================================================

        http_api.add_routes(
            path="/employees",

            methods=[
                apigwv2.HttpMethod.GET,
            ],

            integration=(
                read_employee_integration
            ),
        )

        # =========================================================
        # GET /employees/{id}
        # =========================================================

        http_api.add_routes(
            path="/employees/{id}",

            methods=[
                apigwv2.HttpMethod.GET,
            ],

            integration=(
                read_employee_integration
            ),
        )

        # =========================================================
        # PATCH /employees/{id}
        # =========================================================

        http_api.add_routes(
            path="/employees/{id}",

            methods=[
                apigwv2.HttpMethod.PATCH,
            ],

            integration=(
                update_employee_integration
            ),
        )

        # =========================================================
        # DELETE /employees/{id}
        # =========================================================

        http_api.add_routes(
            path="/employees/{id}",

            methods=[
                apigwv2.HttpMethod.DELETE,
            ],

            integration=(
                delete_employee_integration
            ),
        )

        # =========================================================
        # OUTPUTS
        # =========================================================

        cdk.CfnOutput(
            self,
            "HttpApiUrl",
            value=http_api.api_endpoint,
        )

        cdk.CfnOutput(
            self,
            "CreateUserLambdaArn",
            value=create_user.function_arn,
        )

        cdk.CfnOutput(
            self,
            "ReadUserLambdaArn",
            value=read_user.function_arn,
        )

        cdk.CfnOutput(
            self,
            "CreateEmployeeLambdaArn",
            value=create_employee.function_arn,
        )

        cdk.CfnOutput(
            self,
            "ReadEmployeeLambdaArn",
            value=read_employee.function_arn,
        )

        cdk.CfnOutput(
            self,
            "UpdateEmployeeLambdaArn",
            value=update_employee.function_arn,
        )

        cdk.CfnOutput(
            self,
            "DeleteEmployeeLambdaArn",
            value=delete_employee.function_arn,
        )

        cdk.CfnOutput(
            self,
            "DependencyLayerArn",
            value=(
                dependency_layer.layer_version_arn
            ),
        )