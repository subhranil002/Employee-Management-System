import json

from common.response import error, success
from common.service import (
    EmployeeService,
    InvalidEmployeeError,
)


service = EmployeeService()


def handler(event, context):

    try:
        body = event.get("body")

        if not body:
            return error(
                400,
                "request body is required",
            )

        if isinstance(body, str):
            body = json.loads(body)

        headers = {k.lower(): v for k, v in (event.get("headers") or {}).items()}
        created_by = headers.get("x-user-id")
        
        if not created_by:
            return error(
                400,
                "x-user-id header is required",
            )
            
        body["createdBy"] = created_by

        employee = service.create(body)

        return success(
            201,
            "Employee created successfully",
            employee,
        )

    except json.JSONDecodeError:
        return error(
            400,
            "invalid JSON",
        )

    except InvalidEmployeeError as exc:
        return error(
            400,
            str(exc),
        )

    except Exception as exc:
        print(
            f"Create employee error: {exc}"
        )

        return error(
            500,
            "internal server error",
        )
