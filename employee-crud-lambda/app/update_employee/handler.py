import json

from common.response import error, success
from common.service import (
    EmployeeNotFoundError,
    EmployeeService,
    InvalidEmployeeError,
)

service = EmployeeService()


def handler(event, context):
    try:
        path_parameters = event.get("pathParameters") or {}
        employee_id = path_parameters.get("id")
        if not employee_id:
            return error(400, "employee id is required")

        body = event.get("body")
        if not body:
            return error(400, "request body is required")

        if isinstance(body, str):
            body = json.loads(body)

        # Extract tenant owner user id from request headers
        headers = {k.lower(): v for k, v in (event.get("headers") or {}).items()}
        created_by = headers.get("x-user-id")
        if not created_by:
            return error(400, "x-user-id header is required")

        # Apply partial update to employee document
        employee = service.update(employee_id, body, created_by)
        return success(200, "Employee updated successfully", employee)

    except json.JSONDecodeError:
        return error(400, "invalid JSON")

    except InvalidEmployeeError as exc:
        return error(400, str(exc))

    except EmployeeNotFoundError:
        return error(404, "employee not found")

    except Exception as exc:
        print(f"Update employee error: {exc}")
        return error(500, "internal server error")
